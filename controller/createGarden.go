package controller

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shiibs/go-garden-planner/database"
	"github.com/shiibs/go-garden-planner/model"
)

type Form struct {
    GardenName     string        `json:"gardenName"`
    StartDate      string        `json:"startDate"`
    Rows           string        `json:"rows"`
    Columns        string        `json:"columns"`
    AddedPlantList []model.Plant `json:"addedPlantList"`
}

func PostGardenPlanner(c *fiber.Ctx) error {
    context := fiber.Map{
        "statusText": "OK",
        "message":    "Create garden layout route",
    }
    email := c.Locals("email")

    if email == nil {
        log.Println("Email not found")
        context["msg"] = "Email not found."
        return c.Status(fiber.StatusUnauthorized).JSON(context)
    }

    var user model.User

    if err := database.DBConn.Where("email = ?", email).First(&user).Error; err != nil {
        log.Println("User not found.")
        context["msg"] = "User not found."
        return c.Status(fiber.StatusBadRequest).JSON(context)
    }

    submittedForm := new(Form)

    if err := c.BodyParser(submittedForm); err != nil {
        log.Println("Error in parsing request:", err)
        return respondWithError(c, fiber.StatusInternalServerError, "Error in parsing request")
    }

    // Validate form data
    rows, columns, err := validateFormData(submittedForm)
    if err != nil {
        log.Println("Error in form data:", err)
        return respondWithError(c, fiber.StatusBadRequest, "Invalid form data")
    }

    // Parse start date
    plantingDate, err := time.Parse("2006-01-02", submittedForm.StartDate)
    if err != nil {
        log.Println("Error in parsing Date:", err)
        return respondWithError(c, fiber.StatusBadRequest, "Error in parsing Date")
    }

    // Create garden layout
    garden := ArrangePlants(submittedForm.AddedPlantList, rows, columns)
    monthlyCareDates := generateMonthlyCareDates(plantingDate, 12)

   
    gardenLayout := model.GardenLayout{
        Name:         submittedForm.GardenName,
        StartDate:    plantingDate,
        GardenLayout: garden,
        UserID:       user.ID,
        CareDates:  monthlyCareDates,
    }

    // Save garden layout to database
    if err := saveGardenLayout(&gardenLayout); err != nil {
        log.Println("Error in saving data:", err)
        return respondWithError(c, fiber.StatusInternalServerError, "Error in saving data")
    }

      // Generate and save replanting schedules
      if err := saveReplantingSchedules(submittedForm.AddedPlantList, gardenLayout.ID, plantingDate); err != nil {
        log.Println("Error in saving schedule:", err)
        return respondWithError(c, fiber.StatusInternalServerError, "Error in saving schedule")
    }


    // Send email notification
    if err := SendEmail(user.Email, gardenLayout); err != nil {
        log.Println("Error sending email:", err)
        // We can choose to ignore this error if the main function is successful
        // return respondWithError(c, fiber.StatusInternalServerError, "Error sending email")
    }

    context["gardenId"] = gardenLayout.ID
    return c.Status(fiber.StatusCreated).JSON(context)
}

func respondWithError(c *fiber.Ctx, status int, message string) error {
    return c.Status(status).JSON(fiber.Map{
        "statusText": "Error",
        "message":    message,
    })
}

func validateFormData(form *Form) (int, int, error) {
    rows, err := strconv.Atoi(form.Rows)
    if err != nil || rows <= 0 {
        return 0, 0, err
    }
    columns, err := strconv.Atoi(form.Columns)
    if err != nil || columns <= 0 {
        return 0, 0, err
    }
    return rows, columns, nil
}

func saveGardenLayout(gardenLayout *model.GardenLayout) error {
    result := database.DBConn.Create(gardenLayout)
    return result.Error
}

func saveReplantingSchedules(plants []model.Plant, gardenID uint, startDate time.Time) error {
    for _, plant := range plants {
        replantingDates := generateReplantingDates(plant, 12, startDate)
        for plantName, dates := range replantingDates {
            schedule := model.Schedule{
                PlantName:     plantName,
                GardenID:      gardenID,
                PlantingDates: dates,
            }
            if result := database.DBConn.Create(&schedule); result.Error != nil {
                return result.Error
            }
        }
    }
    return nil
}

func generateMonthlyCareDates(startDate time.Time, numMonths int) []time.Time {
    var careDates []time.Time
    for i := 0; i < numMonths; i++ {
        careDate := startDate.AddDate(0, i, 0)
        careDates = append(careDates, careDate)
    }
    return careDates
}

func generateReplantingDates(plant model.Plant, numMonths int, startDate time.Time) map[string][]time.Time {
    endDate := startDate.AddDate(0, numMonths, 0)
    replantingDates := make(map[string][]time.Time)
    plantingDate := startDate

    for plantingDate.Before(endDate) {
        plantingDate = plantingDate.AddDate(0, 0, plant.ReplantFrequencyDays)
        if plantingDate.Before(endDate) {
            replantingDates[plant.Name] = append(replantingDates[plant.Name], plantingDate)
        }
    }

    return replantingDates
}

func ArrangePlants(plants []model.Plant, rows, cols int) model.Garden {
    garden := model.NewGarden(rows, cols)

    isOuter := func(row, col int) bool {
        return row == 0 || row == rows-1 || col == 0 || col == cols-1
    }

    // First pass: Place plants that need to be on the outer edges
    for i := range plants {
        plant := &plants[i]
        if plant.PlantsPerSquare == 1 {
            for row := 0; row < rows; row++ {
                for col := 0; col < cols; col++ {
                    if isOuter(row, col) && garden[row][col].ID == 0 {
                        garden[row][col] = *plant
                        plant.Count--
                        if plant.Count == 0 {
                            break
                        }
                    }
                }
                if plant.Count == 0 {
                    break
                }
            }
        }
    }

    // Second pass: Place remaining plants in available spaces
    for i := range plants {
        plant := &plants[i]
        if plant.Count > 0 {
            for row := 0; row < rows; row++ {
                for col := 0; col < cols; col++ {
                    if garden[row][col].ID == 0 && !isAdjacentToEnemyPlant(garden, row, col, plant.EnemyPlants) {
                        if row == 0 || !isAdjacentToEnemyPlant(garden, row-1, col, plant.EnemyPlants) {
                            garden[row][col] = *plant
                            plant.Count--
                            if plant.Count == 0 {
                                break
                            }
                        }
                    }
                }
                if plant.Count == 0 {
                    break
                }
            }
        }
    }

    return garden
}

func isAdjacentToEnemyPlant(garden model.Garden, row, col int, enemyPlants []model.Enemy) bool {
    adjacentPositions := [][2]int{
        {row - 1, col},
        {row + 1, col},
        {row, col - 1},
        {row, col + 1},
    }

    enemyPlantNames := make(map[uint]bool)
    for _, enemy := range enemyPlants {
        enemyPlantNames[enemy.EnemyID] = true
    }

    for _, pos := range adjacentPositions {
        r, c := pos[0], pos[1]
        if r >= 0 && r < len(garden) && c >= 0 && c < len(garden[r]) {
            if enemyPlantNames[garden[r][c].ID] {
                return true
            }
        }
    }

    return false
}
// func (g Garden) Print() {
// 	for _, row := range g {
// 		for _, plant := range row {
// 			if plant.Name == "" {
// 				fmt.Print(" - ")
// 			} else {
// 				fmt.Printf("%s ", plant.Name)
// 			}
// 		}
// 		fmt.Println()
// 	}
// }





