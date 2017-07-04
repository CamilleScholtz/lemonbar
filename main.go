package main

import (
	"fmt"
	"log"
)

func main() {
	b, err := newBar()
	if err != nil {
		log.Fatal(err)
	}

	// Run bar blocks.
	go b.workspaceFun()
	go b.memoryFun()
	go b.temperatureFun()
	go b.weatherFun()
	go b.windowFun()
	go b.musicFun()
	go b.todoFun()
	go b.clockFun()

	for {
		// Align to the left side of the bar.
		fmt.Print("%{l}")

		// Workspace icon.
		fmt.Print("%{B#26241C}")
		fmt.Print("%{F#C4C1B0}")
		fmt.Print("  ")

		// Workspace text.
		fmt.Print("%{B#C4C1B0}")
		fmt.Print("%{F#423F31} ")
		fmt.Print(" " + b.workspace + "  ")
		fmt.Print("%{B#423F31} ")

		// Memory icon.
		fmt.Print("%{B#C4C1B0}")
		fmt.Print("  ")

		// Memory text.
		fmt.Print("%{B#687E5A}")
		fmt.Print("%{F#26241C}")
		fmt.Print("  " + b.memory + "  ")
		fmt.Print("%{B#423F31} ")

		// Temperature icon.
		fmt.Print("%{B#C4C1B0}")
		fmt.Print("%{F#423F31} ")
		fmt.Print("  ")

		// Temperature text.
		fmt.Print("%{B#9CB4A6}")
		fmt.Print("  " + b.temperature + "  ")
		fmt.Print("%{B#423F31} ")

		// Window text.
		fmt.Print("%{F#C4C1B0}")
		fmt.Print("   " + b.window + " ")

		// Align to the right side of the bar.
		fmt.Print("%{r}")

		// Music icon.
		fmt.Print("%{B#26241C}")
		fmt.Print("  ")

		// Music text.
		fmt.Print("%{B#C4C1B0}")
		fmt.Print("%{F#423F31}")
		fmt.Print("  " + b.music + "  ")

		// Music state icon.
		if b.musicState {
			fmt.Print("%{B#8E2F34}")
			fmt.Print("%{F#C4C1B0}")
			fmt.Print("  ")
		} else {
			fmt.Print("%{B#687E5A}")
			fmt.Print("%{F#26241C}")
			fmt.Print("  ")
		}
		fmt.Print("%{B#423F31} ")

		// Todo icon.
		fmt.Print("%{B#C4C1B0}")
		fmt.Print("%{F#423F31}")
		fmt.Print("  ")

		// Todo text.
		fmt.Print("%{B#9F8C7C}")
		fmt.Print("%{F#C4C1B0}")
		fmt.Print("  " + b.todo + "  ")
		fmt.Print("%{B#423F31} ")

		// Weather icon
		fmt.Print("%{B#C4C1B0}")
		fmt.Print("%{F#26241C}")
		switch b.weatherState {
		case 0:
			fmt.Print("  ")
		case 1:
			fmt.Print("  ")
		case 2, 3:
			fmt.Print("  ")
		case 4:
			fmt.Print("  ")
		case 5:
			fmt.Print("  ")
		case 6:
			fmt.Print("  ")
		case 7:
			fmt.Print("  ")
		case 8:
			fmt.Print("  ")
		}

		// Weather text
		fmt.Print("%{B#988871}")
		fmt.Print("  " + b.weather + "  ")
		fmt.Print("%{B#423F31} ")

		// Clock icon.
		fmt.Print("%{B#26241C}")
		fmt.Print("%{F#C4C1B0}")
		fmt.Print("  ")

		// Clock text.
		fmt.Print("%{B#C4C1B0}")
		fmt.Print("%{F#423F31}")
		fmt.Print("  " + b.clock + "  ")
		fmt.Print("%{B#423F31}")

		fmt.Println()

		// Wait till there is a block that has updated.
		<-b.Done
	}
}
