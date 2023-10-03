package myparser

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var BeginTime time.Time
var EndTime time.Time

func CheckError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func GenerateTime(timeString string, completeTimeForm bool) time.Time {
	if completeTimeForm {
		t, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			CheckError(err)
		}
		return t
	} else {
		t, err := time.Parse(time.DateOnly, timeString)
		if err != nil {
			CheckError(err)
		}
		return t
	}
}

func ValidateArgBeginTime(args []string) error {
	argString := strings.Join(args, "")
	re := regexp.MustCompile(`^(\d\d\d\d)-(\d\d)-(\d\d)(?:T(\d\d):(\d\d):(\d\d(?:\.\d{1,9})?)(Z|\+\d\d:\d\d|\-\d\d:\d\d))?$`)
	m := re.MatchString(argString)
	if !m {
		return fmt.Errorf("%v is invalid, only accepts RFC3339 format", argString)
	}
	//fmt.Println(re.FindAllString(argString, -1))
	matched := re.FindAllStringSubmatch(argString, -1)
	//fmt.Printf("%q\n", matched)
	completeTimeForm := true
	for _, v := range matched[0] {
		if v == "" {
			completeTimeForm = false
			fmt.Printf("WARN: Incomplete start time detected. The program assumes it is %v.\n", matched[0][0])
			break
		}
	}
	t := GenerateTime(matched[0][0], completeTimeForm)
	BeginTime = t.UTC()
	EndTime = t.Add(time.Hour * 1) // by default generated logs will have one hour duration
	return nil
}

func ValidateArgEndTime(args []string) error {
	argString := strings.Join(args, "")
	if argString == "now" {
		EndTime = time.Now().UTC() // overwrite the default one hour duration
		return nil
	}
	re := regexp.MustCompile(`^(\d\d\d\d)-(\d\d)-(\d\d)(?:T(\d\d):(\d\d):(\d\d(?:\.\d{1,9})?)(Z|\+\d\d:\d\d|\-\d\d:\d\d))?$`)
	m := re.MatchString(argString)
	if !m {
		return fmt.Errorf("%v is invalid, only accepts RFC3339 format", argString)
	}
	//fmt.Println(re.FindAllString(argString, -1))
	matched := re.FindAllStringSubmatch(argString, -1)
	//fmt.Printf("%q\n", matched)
	completeTimeForm := true
	for _, v := range matched[0] {
		if v == "" {
			completeTimeForm = false
			fmt.Printf("WARN: Incomplete end time detected. The program assumes it is %v.\n", matched[0][0])
			break
		}
	}
	t := GenerateTime(matched[0][0], completeTimeForm)

	diff := t.Sub(BeginTime)
	if diff < 0 {
		return fmt.Errorf("%v is invalid, it should be later than the begin time", argString)
	}

	EndTime = t.UTC() // overwrite the default one hour duration
	return nil
}

func GenerateArgTimeExamples(IncludeNow bool) string {
	message := "Example 1: 2023-02-02 (will be assumed 2023-02-02T00:00:00Z)\n" +
		"Example 2: 2023-02-02T09:40:00Z\n" +
		"Example 3: 2023-02-02T09:40:00+08:00\n" +
		"Example 4: 2023-02-02T09:40:00-12:00\n"

	if IncludeNow {
		message += "Example 5: now (will generate the current time)\n"
	}

	return message
}

func ValidateArgDuration(args []string) error {
	argString := strings.Join(args, "")
	re := regexp.MustCompile(`^([0-9]+)(s|m|h|d)$`)
	m := re.MatchString(argString)
	if !m {
		return fmt.Errorf("%v is invalid, only accepts N(s|m|h|d), e.g. 7d for 7 days", argString)
	}
	matched := re.FindAllStringSubmatch(argString, -1)
	switch matched[0][2] {
	case "s":
		i, _ := strconv.ParseInt(matched[0][1], 10, 64)
		EndTime = BeginTime.Add(time.Second * time.Duration(i))
	case "m":
		i, _ := strconv.ParseInt(matched[0][1], 10, 64)
		EndTime = BeginTime.Add(time.Minute * time.Duration(i))
	case "h":
		i, _ := strconv.ParseInt(matched[0][1], 10, 64)
		EndTime = BeginTime.Add(time.Hour * time.Duration(i))
	case "d":
		i, _ := strconv.ParseInt(matched[0][1], 10, 64)
		EndTime = BeginTime.Add(time.Hour * 24 * time.Duration(i))
	default:
		log.Fatalf("Invalid duration input, only accepts N(s|m|h|d), e.g. 7d for 7 days")
	}
	return nil
}
