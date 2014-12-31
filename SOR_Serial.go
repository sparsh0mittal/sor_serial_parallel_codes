//The use of this file is subject to the license terms in the LICENSE file 
//found in the top-level directory of this distribution. No part of this file,
// may be copied, modified, propagated, or distributed except according to 
//the terms contained in the LICENSE file.
//Author Sparsh Mittal sparsh0mittal@gmail.com 

package main




import "fmt"

//import "sync"
import "io"
import "os"

//import "strings"
import "time"
import "strconv"
import "bytes"

const gridSize int = 4096
const INITIAL_GUESS float64 = 0.0
const MAXSTEPS int = 500000        /* Maximum number of iterations               */
const TOL_VAL float64 = 0.00001   /* Numerical Tolerance */
const PI_VAL float64 = 3.14159265 /* pi */
const NUM_CHECK int = 4000
const omega float64 = 0.376
const one_minus_omega float64 = 1- omega 

//const numberOfSlaves int = 2;

//var wg      sync.WaitGroup

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

func main() {
	fmt.Println("GridSize ",gridSize, " MaxStep ", MAXSTEPS, " Accuracy ", TOL_VAL )
	var outputFileName bytes.Buffer
	var maxError float64 = 0.0
	var isConverged bool = false

	outputFileName.WriteString("Serial")
	outputFileName.WriteString(strconv.Itoa(gridSize))

	var timeFileName bytes.Buffer
	timeFileName.WriteString("Time_Serial_")
	timeFileName.WriteString(strconv.Itoa(gridSize))
	
	
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
	var sum float64 = 0.0

	startTime := time.Now()
	for iter := 1; iter <= MAXSTEPS; iter++ {
		maxError = 0.0
		for i := 1; i < gridSize+1; i++ {
			for j := 1; j < gridSize+1; j++ {
				if (i+j)%2 == 0 {
					sum = (gridInfo[i][j+1] + gridInfo[i+1][j] + gridInfo[i-1][j] + gridInfo[i][j-1]) * 0.25

					maxError = MAX_FUNC(ABS_VAL(omega*(sum-gridInfo[i][j])), maxError)
					gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum
				}

			}
		}

		for i := 1; i < gridSize+1; i++ {
			for j := 1; j < gridSize+1; j++ {
				if (i+j)%2 == 1 {
					sum = (gridInfo[i][j+1] + gridInfo[i+1][j] + gridInfo[i-1][j] + gridInfo[i][j-1]) * 0.25

					maxError = MAX_FUNC(ABS_VAL(omega*(sum-gridInfo[i][j])), maxError)
					gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum
				}

			}
		}

		if iter%NUM_CHECK == 0 {
			fmt.Println("Iter ", iter, " Error ", maxError)
			if maxError < TOL_VAL {
				isConverged = true
				break
			} 
		}
	}
	stopTime := time.Now()
	fmt.Println("isConverged", isConverged)
	timeTaken := stopTime.Sub(startTime)

	fmt.Println("Total time was ", timeTaken.Seconds())
	timestring := strconv.FormatInt( int64 (timeTaken.Seconds()), 10)
	io.WriteString(ftime, timestring)
	spacestring := " "

	for i := 0; i < gridSize+2; i++ {
		for j := 0; j < gridSize+2; j++ {
			
			s := strconv.FormatFloat(gridInfo[i][j], 'e', 6, 64)			
			io.WriteString(foutput,  s)
			io.WriteString(foutput, spacestring)
		}
		
		s2 := "\n"
		io.WriteString(foutput, s2);
	}
}
