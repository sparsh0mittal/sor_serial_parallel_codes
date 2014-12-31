//The use of this file is subject to the license terms in the LICENSE file 
//found in the top-level directory of this distribution. No part of this file,
// may be copied, modified, propagated, or distributed except according to 
//the terms contained in the LICENSE file.
//Author Sparsh Mittal sparsh0mittal@gmail.com


package main

import "fmt"

import "sync"
import "io"
import "os"
import "runtime"

//import "strings"
import "time"
import "strconv"
import "bytes"

const gridSize int = 4096
const INITIAL_GUESS float64 = 0.0
const MAXSTEPS int = 500000       /* Maximum number of iterations               */
const TOL_VAL float64 = 0.00001   /* Numerical Tolerance */
const PI_VAL float64 = 3.14159265 /* pi */
const NUM_CHECK int = 4000
const omega float64 = 0.376
const one_minus_omega float64 = 1 - omega

var channelLock = make(chan int, 1)

const numberOfSlaves int = 4
const factorValue int = gridSize / numberOfSlaves

var wg sync.WaitGroup

var gridInfo = make([][]float64, gridSize+2)

func MAX_FUNC(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func ABS_VAL(a float64) float64 {

	if a > 0 {
		return a
	}
	return -a
}

var maxError float64 = 0.0

func main() {

	runtime.GOMAXPROCS(numberOfSlaves)
	fmt.Println("GridSize ", gridSize, " MaxStep ", MAXSTEPS, " Accuracy ", TOL_VAL)

	if gridSize%numberOfSlaves != 0 {
		panic("not integer")
	}

	var outputFileName bytes.Buffer

	var isConverged bool = false

	outputFileName.WriteString("Parallel_")
	outputFileName.WriteString(strconv.Itoa(gridSize))
	outputFileName.WriteString("_On")
	outputFileName.WriteString(strconv.Itoa(numberOfSlaves))

	var timeFileName bytes.Buffer
	timeFileName.WriteString("Time_Parallel_")
	timeFileName.WriteString(strconv.Itoa(gridSize))
	timeFileName.WriteString("_On")
	timeFileName.WriteString(strconv.Itoa(numberOfSlaves))

	ftime, err2 := os.Create(timeFileName.String())
	if err2 != nil {
		panic(err2)
	}
	defer ftime.Close()

	foutput, err1 := os.Create(outputFileName.String())
	if err1 != nil {
		panic(err1)
	}
	defer foutput.Close()

	// allocate composed 2d array

	for i := range gridInfo {
		gridInfo[i] = make([]float64, gridSize+2)
	}

	for i := 0; i < gridSize+2; i++ {
		for j := 0; j < gridSize+2; j++ {
			if i == 0 {
				gridInfo[i][j] = 1.0
			} else {
				gridInfo[i][j] = INITIAL_GUESS
			}

		}
	}

	var shouldCheck bool = false
	startTime := time.Now()

	channelLock <- 1
	for iter := 1; iter <= MAXSTEPS; iter++ {

		shouldCheck = false
		if iter%NUM_CHECK == 0 {
			shouldCheck = true
			maxError = 0.0
		}

		wg.Add(numberOfSlaves)
		for islave := 0; islave < numberOfSlaves; islave++ {
			go sorSolverSlave(islave, 0, shouldCheck)
		}
		wg.Wait()

		wg.Add(numberOfSlaves)
		for islave := 0; islave < numberOfSlaves; islave++ {
			go sorSolverSlave(islave, 1, shouldCheck)
		}
		wg.Wait()

		if shouldCheck {
			fmt.Println("Iter ", iter, " Error ", maxError)
			if maxError < TOL_VAL {
				isConverged = true
				break
			}
		}
	}

	//wg.Wait()
	stopTime := time.Now()
	fmt.Println("isConverged", isConverged)
	timeTaken := stopTime.Sub(startTime)

	fmt.Println("Total time was ", timeTaken.Seconds())

	timestring := strconv.FormatInt(int64(timeTaken.Seconds()), 10)
	io.WriteString(ftime, timestring)

	

}

func sorSolverSlave(myNumber int, whichPhase int, shouldCheckHere bool) {

	//var remainder int =whichPhase;

	var sum float64 = 0

	var iStart int = (myNumber * factorValue) + 1
	var iEnd int = ((myNumber + 1) * factorValue)

	if whichPhase == 0 {
		//i even, j even
		for i := iStart + 1; i <= iEnd; i += 2 {
			for j := 2; j < gridSize+1; j += 2 {

				sum = (gridInfo[i][j+1] + gridInfo[i+1][j] + gridInfo[i-1][j] + gridInfo[i][j-1]) * 0.25

				if shouldCheckHere {
					var errorHere float64 = ABS_VAL(omega * (sum - gridInfo[i][j]))
					<-channelLock
					maxError = MAX_FUNC(errorHere, maxError)
					channelLock <- 1
				}
				gridInfo[i][j] = (one_minus_omega)*gridInfo[i][j] + omega*sum

			}
		}
		//i odd, j odd
		for i := iStart; i <= iEnd; i += 2 {
			for j := 1; j < gridSize+1; j += 2 {

				sum = (gridInfo[i][j+1] + gridInfo[i+1][j] + gridInfo[i-1][j] + gridInfo[i][j-1]) * 0.25

				if shouldCheckHere {
					var errorHere float64 = ABS_VAL(omega * (sum - gridInfo[i][j]))
					<-channelLock
					maxError = MAX_FUNC(errorHere, maxError)
					channelLock <- 1
				}
				gridInfo[i][j] = (one_minus_omega)*gridInfo[i][j] + omega*sum

			}
		}
	} else {
		//i even, j odd
		for i := iStart + 1; i <= iEnd; i += 2 {
			for j := 1; j < gridSize+1; j += 2 {

				sum = (gridInfo[i][j+1] + gridInfo[i+1][j] + gridInfo[i-1][j] + gridInfo[i][j-1]) * 0.25

				if shouldCheckHere {
					var errorHere float64 = ABS_VAL(omega * (sum - gridInfo[i][j]))
					<-channelLock
					maxError = MAX_FUNC(errorHere, maxError)
					channelLock <- 1
				}
				gridInfo[i][j] = (one_minus_omega)*gridInfo[i][j] + omega*sum

			}
		}

		// i odd, j even
		for i := iStart; i <= iEnd; i += 2 {
			for j := 2; j < gridSize+1; j += 2 {

				sum = (gridInfo[i][j+1] + gridInfo[i+1][j] + gridInfo[i-1][j] + gridInfo[i][j-1]) * 0.25

				if shouldCheckHere {
					var errorHere float64 = ABS_VAL(omega * (sum - gridInfo[i][j]))
					<-channelLock
					maxError = MAX_FUNC(errorHere, maxError)
					channelLock <- 1
				}
				gridInfo[i][j] = (one_minus_omega)*gridInfo[i][j] + omega*sum

			}
		}

	}

	wg.Done()
}

