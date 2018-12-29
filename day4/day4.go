package main

import (
	"fmt"
	"github.com/Piszmog/advent-of-code-2018/utils"
	"log"
	"sort"
	"strings"
	"time"
)

const filename = "day4/guardSchedule.txt"

type Guard struct {
	id                 int
	totalMinutesAsleep int
	minutesAsleep      []MinuteRange
}

type MinuteRange struct {
	start int
	end   int
}

type Schedule struct {
	timestamp time.Time
	text      string
}

func main() {
	defer utils.Runtime(time.Now())
	file := utils.OpenFile(filename)
	defer utils.CloseFile(file)
	lines := make(chan string, 200)
	done := make(chan bool)
	go mapSchedule(lines, done)
	utils.ReadTXTFile(file, lines)
	<-done
}

func mapSchedule(lines chan string, done chan bool) {
	var schedules []Schedule
	for line := range lines {
		timestamp := line[1:17]
		scheduleTime, err := time.Parse("2006-01-02 15:04", timestamp)
		if err != nil {
			log.Fatal(err)
		}
		text := line[18:]
		schedules = append(schedules, Schedule{
			timestamp: scheduleTime,
			text:      text,
		})
	}
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].timestamp.Before(schedules[j].timestamp)
	})
	guardId := 0
	guards := make(map[int]Guard)
	var timeAsleep int
	var minuteRange []MinuteRange
	var sleepTime time.Time
	for _, schedule := range schedules {
		if strings.Contains(schedule.text, "Guard") {
			if guardId != 0 {
				guard := guards[guardId]
				if guard.id != 0 {
					guard.totalMinutesAsleep = guard.totalMinutesAsleep + timeAsleep
					guard.minutesAsleep = append(guard.minutesAsleep, minuteRange...)
				} else {
					guard = Guard{
						id:                 guardId,
						totalMinutesAsleep: timeAsleep,
						minutesAsleep:      minuteRange,
					}
				}
				guards[guardId] = guard
				timeAsleep = 0
				minuteRange = minuteRange[:0]
			}
			num, err := fmt.Sscanf(schedule.text, " Guard #%d begins shift", &guardId)
			if err != nil {
				log.Fatal(err)
			}
			if num != 1 {
				log.Fatal("could not parse guard id")
			}
		} else if strings.Contains(schedule.text, "falls") {
			sleepTime = schedule.timestamp
		} else {
			wakeUpTime := schedule.timestamp //.Add(time.Duration(-1) * time.Minute)
			minuteRange = append(minuteRange, MinuteRange{
				start: sleepTime.Minute(),
				end:   wakeUpTime.Minute(),
			})
			asleepMinutes := wakeUpTime.Minute() - sleepTime.Minute()
			timeAsleep = timeAsleep + asleepMinutes
		}
	}
	guardAsleepLongest := 0
	timeAsleepLongest := 0
	for id, guard := range guards {
		if timeAsleepLongest < guard.totalMinutesAsleep {
			guardAsleepLongest = id
			timeAsleepLongest = guard.totalMinutesAsleep
		}
	}
	guard := guards[guardAsleepLongest]
	minuteAsleepTimes := make(map[int]int)
	for _, minuteRange := range guard.minutesAsleep {
		start := minuteRange.start
		end := minuteRange.end
		for i := start; i <= end; i++ {
			minuteAsleepTimes[start]++
		}
	}
	minuteMostAsleep := 0
	mostNumberOfTimesAsleep := 0
	for minute, numberOfTimesAsleep := range minuteAsleepTimes {
		if mostNumberOfTimesAsleep < numberOfTimesAsleep {
			mostNumberOfTimesAsleep = numberOfTimesAsleep
			minuteMostAsleep = minute
		}
	}
	fmt.Printf("Guard %d slept the longest for %d minutes. Minute most asleep is %d\n", guardAsleepLongest, timeAsleepLongest, minuteMostAsleep)
	close(done)
}
