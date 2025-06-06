package tools

import (
	"github.com/gocolly/colly/v2"
	"regexp"
	"strings"
)

type meal struct {
	Name           string
	Allergens      string
	Additives      string
	PriceStudents  string
	PriceEmployees string
	PriceGuests    string
}

type Meals struct {
	Monday    []meal
	Tuesday   []meal
	Wednesday []meal
	Thursday  []meal
	Friday    []meal
	Saturday  []meal
}

func PrepareMessageForMeals(meals Meals, detailed bool, days ...string) string {
	if (days == nil) || (len(days) == 0) {
		days = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	}

	message := "Mensa Furtwangen\n\n"
	if StringInSlice("Monday", days) {
		message += "Montag\n"
		for _, meal := range meals.Monday {
			message += prepareMessageDay(meal, detailed)
		}
	}
	if StringInSlice("Tuesday", days) {
		message += "Dienstag\n"
		for _, meal := range meals.Tuesday {
			message += prepareMessageDay(meal, detailed)
		}
	}
	if StringInSlice("Wednesday", days) {
		message += "Mittwoch\n"
		for _, meal := range meals.Wednesday {
			message += prepareMessageDay(meal, detailed)
		}
	}
	if StringInSlice("Thursday", days) {
		message += "Donnerstag\n"
		for _, meal := range meals.Thursday {
			message += prepareMessageDay(meal, detailed)
		}
	}
	if StringInSlice("Friday", days) {
		message += "Freitag\n"
		for _, meal := range meals.Friday {
			message += prepareMessageDay(meal, detailed)
		}
	}
	if StringInSlice("Saturday", days) {
		message += "Samstag\n"
		for _, meal := range meals.Saturday {
			message += prepareMessageDay(meal, detailed)
		}
	}

	return message
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func prepareMessageDay(meal meal, detailed bool) string {
	message := "    " + meal.Name + "\n"
	if detailed {
		message += "    Allergene: " + meal.Allergens + "\n"
		message += "    Zusatzstoffe: " + meal.Additives + "\n"
		message += "    Preis Studierende: " + meal.PriceStudents + "\n"
		message += "    Preis Beschäftigte: " + meal.PriceEmployees + "\n"
		message += "    Preis Gäste: " + meal.PriceGuests + "\n"
	} else {
		message += "    Preis Studierende: " + meal.PriceStudents + "\n"
	}
	message += "\n"
	return message
}

func LoadMeals() Meals {

	allMeals := Meals{}

	c := colly.NewCollector()
	// Monday
	// OnHTML is called for each element found with the specified selector
	// in this case, the selector is a div with the id "tab-mon" and a child div with the class "col-span-1"
	c.OnHTML("div#tab-mon div.col-span-1", func(e *colly.HTMLElement) {
		// get text
		text := cleanText(e.Text)
		// parse text into meal
		myMeal := parseMeal(text)
		//add meal to corresponding day
		allMeals.Monday = append(allMeals.Monday, myMeal)
	})
	// Tuesday
	c.OnHTML("div#tab-tue div.col-span-1", func(e *colly.HTMLElement) {
		text := cleanText(e.Text)
		myMeal := parseMeal(text)
		allMeals.Tuesday = append(allMeals.Tuesday, myMeal)
	})
	// Wednesday
	c.OnHTML("div#tab-wed div.col-span-1", func(e *colly.HTMLElement) {
		text := cleanText(e.Text)
		myMeal := parseMeal(text)
		allMeals.Wednesday = append(allMeals.Wednesday, myMeal)
	})
	// Thursday
	c.OnHTML("div#tab-thu div.col-span-1", func(e *colly.HTMLElement) {
		text := cleanText(e.Text)
		myMeal := parseMeal(text)
		allMeals.Thursday = append(allMeals.Thursday, myMeal)
	})
	// Friday
	c.OnHTML("div#tab-fri div.col-span-1", func(e *colly.HTMLElement) {
		text := cleanText(e.Text)
		myMeal := parseMeal(text)
		allMeals.Friday = append(allMeals.Friday, myMeal)
	})
	// Saturday
	c.OnHTML("div#tab-sat div.col-span-1", func(e *colly.HTMLElement) {
		text := cleanText(e.Text)
		myMeal := parseMeal(text)
		allMeals.Saturday = append(allMeals.Saturday, myMeal)
	})

	err := c.Visit("https://www.swfr.de/essen/mensen-cafes-speiseplaene/mensa-furtwangen")
	if err != nil {
		return Meals{}
	}
	return allMeals
}

func cleanText(text string) string {
	// remove all newlines and tabs, replace them with spaces
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	// replace multiple spaces with one space
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	return text
}

func parseMeal(text string) meal {
	myMeal := meal{}
	chunks := strings.Split(text, "enthält Allergene: ")
	if len(chunks) == 2 {
		myMeal.Allergens = chunks[1]
	} else if len(chunks) == 1 {
		// triggers if the string "enthält Allergene: " is not found
		myMeal.Allergens = "Keine"
	} else {
		// something went completely wrong - it should not happen
		myMeal.Allergens = "Error on Allergenes - Hoping for the best"
	}
	chunks = strings.Split(chunks[0], "Kennzeichnungen/Zusatzstoffe: ")
	if len(chunks) == 2 {
		myMeal.Additives = chunks[1]
	} else if len(chunks) == 1 {
		myMeal.Additives = "Keine"
	} else {
		myMeal.Additives = "Error on Additives - Hoping for the best"
	}
	chunks = strings.Split(chunks[0], "Gäste")
	if len(chunks) == 2 {
		myMeal.PriceGuests = chunks[1]
	} else if len(chunks) == 1 {
		myMeal.PriceGuests = "Keine Angabe"
	} else {
		myMeal.PriceGuests = "Error on PriceGuests - Hoping for the best"
	}
	chunks = strings.Split(chunks[0], "Beschäftigte")
	if len(chunks) == 2 {
		myMeal.PriceEmployees = chunks[1]
	} else if len(chunks) == 1 {
		myMeal.PriceEmployees = "Keine Angabe"
	} else {
		myMeal.PriceEmployees = "Error on PriceEmployees - Hoping for the best"
	}
	chunks = strings.Split(chunks[0], "Studierende, Schüler")
	if len(chunks) == 2 {
		myMeal.PriceStudents = chunks[1]
	} else if len(chunks) == 1 {
		myMeal.PriceStudents = "Keine Angabe"
	} else {
		myMeal.PriceStudents = "Error on PriceStudents - Hoping for the best"
	}
	myMeal.Name = chunks[0]
	myMeal.Name = regexp.MustCompile(`Essen \d `).ReplaceAllString(myMeal.Name, " ")
	myMeal.Name = regexp.MustCompile(`Preise \+ Kennzeichnungen`).ReplaceAllString(myMeal.Name, " ")

	// convert camelCase to "normal case" and remove leading and trailing whitespaces
	myMeal.Name, _ = camelCaseToNormalCase(strings.TrimSpace(myMeal.Name))
	myMeal.Allergens, _ = camelCaseToNormalCase(strings.TrimSpace(myMeal.Allergens))
	myMeal.Additives, _ = camelCaseToNormalCase(strings.TrimSpace(myMeal.Additives))
	myMeal.PriceStudents, _ = camelCaseToNormalCase(strings.TrimSpace(myMeal.PriceStudents))
	myMeal.PriceEmployees, _ = camelCaseToNormalCase(strings.TrimSpace(myMeal.PriceEmployees))
	myMeal.PriceGuests, _ = camelCaseToNormalCase(strings.TrimSpace(myMeal.PriceGuests))

	return myMeal
}

func camelCaseToNormalCase(inputCamelCaseStr string) (string, int) {
	// standard regex for this is `[A-Z][^A-Z]*`
	// I added a second part to also match prices like "2.50 €" that would be lost otherwise
	re := regexp.MustCompile(`([A-z][^A-Z]*|\d?.?\d+ ?€?)`)
	parts := re.FindAllString(inputCamelCaseStr, -1)
	for index := range parts {
		parts[index] = strings.ToLower(parts[index])
	}
	return strings.Join(parts, " "), len(parts)
}
