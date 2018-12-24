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
	id            int
	timeAsleep    []time.Duration
	minutesAsleep []MinuteRange
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
	var timeAsleep []time.Duration
	var minuteRange []MinuteRange
	var sleepTime time.Time
	for _, schedule := range schedules {
		if strings.Contains(schedule.text, "Guard") {
			if guardId != 0 {
				guard := guards[guardId]
				if guard.id != 0 {
					guard.timeAsleep = append(guard.timeAsleep, timeAsleep...)
					guard.minutesAsleep = append(guard.minutesAsleep, minuteRange...)
				} else {
					guard = Guard{
						id:            guardId,
						timeAsleep:    timeAsleep,
						minutesAsleep: minuteRange,
					}
				}
				guards[guardId] = guard
				timeAsleep = timeAsleep[:0]
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
			asleepDuration := wakeUpTime.Sub(sleepTime)
			timeAsleep = append(timeAsleep, asleepDuration)
		}
	}
	var longestTotalAsleep time.Duration
	var longestAsleepGuard int
	for id, guard := range guards {
		var totalTimeAsleep time.Duration
		var longestAsleepDuration time.Duration
		for _, asleepDuration := range guard.timeAsleep {
			if longestAsleepDuration < asleepDuration {
				longestAsleepDuration = asleepDuration
			}
			totalTimeAsleep = totalTimeAsleep + asleepDuration
		}
		if longestTotalAsleep < totalTimeAsleep {
			longestTotalAsleep = totalTimeAsleep
			longestAsleepGuard = id
		}
	}
	guard := guards[longestAsleepGuard]
	minuteAsleepOccurrences := make(map[int]int)
	for _, minuteRange := range guard.minutesAsleep {
		for i := minuteRange.start; i <= minuteRange.end; i++ {
			minuteAsleepOccurrences[i]++
		}
	}
	var minuteMostAsleep int
	var mostTimeAsleep int
	for minute, numberOfTimes := range minuteAsleepOccurrences {
		if mostTimeAsleep < numberOfTimes {
			minuteMostAsleep = minute
			mostTimeAsleep = numberOfTimes
		}
	}
	fmt.Printf("Guard %d slept the longest for %d minutes\n", longestAsleepGuard, minuteMostAsleep)
	close(done)
}
