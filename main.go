package main

/*02/17/20223:
def-x=250
def-y=200
dl5=600
*/

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type dataPoint struct {
	x float64
	y float64
}

var cupDataDeflected [][]dataPoint
var cupDataUndeflected [][]dataPoint
var faceDataDeflected [][]dataPoint
var faceDataUndeflected [][]dataPoint

func main() {

	//program flags (-h for help)
	cup := flag.Bool("cup", false, "use cup data")
	run := flag.Int("run", 1, "which run to use")
	raw := flag.Bool("raw", false, "show raw data")
	n := flag.Int("n", 10, "clustering value for time average")
	flag.Parse()

	//initialize plot
	p := plot.New()

	//sets title based on if we want cup or faceplate data
	if *cup {
		p.Title.Text = "Cup Run " + strconv.Itoa(*run)
	} else {
		p.Title.Text = "Faceplate Run " + strconv.Itoa(*run)
	}

	//axis labels
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Voltage (mV)"

	//function to read files from data folder
	readFile()

	//add points based on if we want raw data or time average
	if *raw {
		err := plotutil.AddLinePoints(p,
			"Deflected", getPoints(*cup, true, *run-1),
			"Undeflected", getPoints(*cup, false, *run-1))
		if err != nil {
			panic(err)
		}
	} else {
		err := plotutil.AddLinePoints(p,
			"Deflected", getTimeAvgPoints(*cup, true, *run-1, *n),
			"Undeflected", getTimeAvgPoints(*cup, false, *run-1, *n))
		if err != nil {
			panic(err)
		}
	}

	// Save the plot to a PNG file.
	if err := p.Save(12*vg.Inch, 9*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

	getMeanSigma(*cup, *run)
}

func getMeanSigma(cup bool, run int) {

	if cup {

		sum := 0.0
		for _, i := range cupDataDeflected[run-1] {
			sum += i.y
		}
		mean := sum / float64(len(cupDataDeflected[run-1]))
		diffSq := 0.0
		for _, i := range cupDataDeflected[run-1] {
			diffSq += math.Pow(i.y-mean, 2)
		}
		sigma := math.Sqrt(diffSq / float64(len(cupDataDeflected[run-1])))
		fmt.Printf("Deflected Mean: %f\nDefelcted Sigma %f\n", mean, sigma)

		sum = 0.0
		for _, i := range cupDataUndeflected[run-1] {
			sum += i.y
		}
		mean = sum / float64(len(cupDataUndeflected[run-1]))
		diffSq = 0.0
		for _, i := range cupDataUndeflected[run-1] {
			diffSq += math.Pow(i.y-mean, 2)
		}
		sigma = math.Sqrt(diffSq / float64(len(cupDataUndeflected[run-1])))
		fmt.Printf("Undeflected Mean: %f\nUndefelcted Sigma %f\n", mean, sigma)
	} else {
		sum := 0.0
		for _, i := range faceDataDeflected[run-1] {
			sum += i.y
		}
		mean := sum / float64(len(faceDataDeflected[run-1]))
		diffSq := 0.0
		for _, i := range faceDataDeflected[run-1] {
			diffSq += math.Pow(i.y-mean, 2)
		}
		sigma := math.Sqrt(diffSq / float64(len(faceDataDeflected[run-1])))
		fmt.Printf("Deflected Mean: %f\nDefelcted Sigma %f\n", mean, sigma)

		sum = 0.0
		for _, i := range faceDataUndeflected[run-1] {
			sum += i.y
		}
		mean = sum / float64(len(faceDataUndeflected[run-1]))
		diffSq = 0.0
		for _, i := range faceDataUndeflected[run-1] {
			diffSq += math.Pow(i.y-mean, 2)
		}
		sigma = math.Sqrt(diffSq / float64(len(faceDataUndeflected[run-1])))
		fmt.Printf("Undeflected Mean: %f\nUndefelcted Sigma %f\n", mean, sigma)
	}

}

// creates a plotter.XYs for the data with a time average of n data points
func getTimeAvgPoints(cup bool, def bool, run int, n int) plotter.XYs {

	//cup or faceplate data
	if cup {
		//deflected on undelfected data
		if def {
			//initalize plotter.XYs with data_length/n points
			pts := make(plotter.XYs, len(cupDataDeflected[run])/n)
			//index over all points
			for i := range pts {
				ysum := 0.0
				//sum up n data values for this points
				for j := 0; j < n; j++ {
					ysum += cupDataDeflected[run][i*n+j].y
				}
				mean := ysum / float64(n)
				diffsq := 0.0
				for j := 0; j < n; j++ {
					diffsq += math.Pow(mean-cupDataDeflected[run][i*n+j].y, 2)
				}
				//sigma := math.Sqrt(diffsq / float64(n))
				//fmt.Printf("Mean: %f Sigma: %f\n", mean, sigma)

				//add average to data point and center at middle time of averaged data points
				pts[i].X = cupDataDeflected[run][i+(n/2)].x
				pts[i].Y = ysum / float64(n)
			}
			return pts
		}
		//see above now with undeflected
		pts := make(plotter.XYs, len(cupDataUndeflected[run])/n)
		for i := range pts {
			ysum := 0.0
			for j := 0; j < n; j++ {
				ysum += cupDataUndeflected[run][i*n+j].y
			}
			mean := ysum / float64(n)
			diffsq := 0.0
			for j := 0; j < n; j++ {
				diffsq += math.Pow(mean-cupDataUndeflected[run][i*n+j].y, 2)
			}
			//sigma := math.Sqrt(diffsq / float64(n))
			//fmt.Printf("Mean: %f Sigma: %f\n", mean, sigma)
			pts[i].X = cupDataUndeflected[run][i+(n/2)].x
			pts[i].Y = ysum / float64(n)
		}
		return pts
	}

	//see above now with faceplate
	if def {
		pts := make(plotter.XYs, len(faceDataDeflected[run])/n)
		for i := range pts {
			ysum := 0.0
			for j := 0; j < n; j++ {
				ysum += faceDataDeflected[run][i*n+j].y
			}
			mean := ysum / float64(n)
			diffsq := 0.0
			for j := 0; j < n; j++ {
				diffsq += math.Pow(mean-faceDataDeflected[run][i*n+j].y, 2)
			}
			//sigma := math.Sqrt(diffsq / float64(n))
			//fmt.Printf("Mean: %f Sigma: %f\n", mean, sigma)
			pts[i].X = faceDataDeflected[run][i+(n/2)].x
			pts[i].Y = ysum / float64(n)
		}
		return pts
	}
	pts := make(plotter.XYs, len(faceDataUndeflected[run])/n)
	for i := range pts {
		ysum := 0.0
		for j := 0; j < n; j++ {
			ysum += faceDataUndeflected[run][i*n+j].y
		}
		mean := ysum / float64(n)
		diffsq := 0.0
		for j := 0; j < n; j++ {
			diffsq += math.Pow(mean-faceDataUndeflected[run][i*n+j].y, 2)
		}
		//sigma := math.Sqrt(diffsq / float64(n))
		//fmt.Printf("Mean: %f Sigma: %f\n", mean, sigma)
		pts[i].X = faceDataUndeflected[run][i+(n/2)].x
		pts[i].Y = ysum / float64(n)
	}
	return pts

}

// get gonumplot points from data
func getPoints(cup bool, def bool, run int) plotter.XYs {

	//see getTimeAvgPoints func *this one is simpler*
	if cup {
		if def {
			pts := make(plotter.XYs, len(cupDataDeflected[run]))
			for i := range pts {
				pts[i].X = cupDataDeflected[run][i].x
				pts[i].Y = cupDataDeflected[run][i].y
			}
			return pts
		}
		pts := make(plotter.XYs, len(cupDataUndeflected[run]))
		for i := range pts {
			pts[i].X = cupDataUndeflected[run][i].x
			pts[i].Y = cupDataUndeflected[run][i].y
		}
		return pts
	}

	if def {
		pts := make(plotter.XYs, len(faceDataDeflected[run]))
		for i := range pts {
			pts[i].X = faceDataDeflected[run][i].x
			pts[i].Y = faceDataDeflected[run][i].y
		}
		return pts
	}
	pts := make(plotter.XYs, len(faceDataUndeflected[run]))
	for i := range pts {
		pts[i].X = faceDataUndeflected[run][i].x
		pts[i].Y = faceDataUndeflected[run][i].y
	}
	return pts

}

// read files and make list
func readFile() {

	//get all files in data directory
	files, err := os.ReadDir("data")

	if err != nil {
		log.Fatal(err)
	}

	//for every file
	for _, file := range files {

		//get the file of interest
		f, err := os.Open("data/" + file.Name())

		if err != nil {
			log.Fatal(err)
		}

		//the data to be added to the list
		var tmpArr []dataPoint

		defer f.Close()

		scanner := bufio.NewScanner(f)

		scanner.Scan()

		var count float64 = 0

		//scan over each line indivivually
		for scanner.Scan() {

			//split and get value as float64
			y, err := strconv.ParseFloat(strings.Split(scanner.Text(), "	")[1], 64)
			if err != nil {
				log.Fatal(err)
			}

			//create a dataPoint struct with data and increase count(time)
			tmp := dataPoint{count, y}
			count++

			//append to tmpArr
			tmpArr = append(tmpArr, tmp)

		}

		//double switch to look at file name and choose which array to append to
		switch strings.Split(file.Name(), "_")[0] {
		case "cup":
			switch strings.Split(file.Name(), "_")[1] {
			case "deflected":
				cupDataDeflected = append(cupDataDeflected, tmpArr)
			case "undeflected":
				cupDataUndeflected = append(cupDataUndeflected, tmpArr)
			}
		case "faceplate":
			switch strings.Split(file.Name(), "_")[1] {
			case "deflected":
				faceDataDeflected = append(faceDataDeflected, tmpArr)
			case "undeflected":
				faceDataUndeflected = append(faceDataUndeflected, tmpArr)
			}
		}

	}

}
