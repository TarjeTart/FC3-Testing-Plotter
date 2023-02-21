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

var deflected []dataPoint
var undeflected []dataPoint

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
	readFile(*cup, *run)

	//add points based on if we want raw data or time average
	if *raw {
		err := plotutil.AddLinePoints(p,
			"Deflected", getPoints(true),
			"Undeflected", getPoints(false))
		if err != nil {
			panic(err)
		}
	} else {
		err := plotutil.AddLinePoints(p,
			"Deflected", getTimeAvgPoints(true, *n),
			"Undeflected", getTimeAvgPoints(false, *n))
		if err != nil {
			panic(err)
		}
	}

	// Save the plot to a PNG file.
	if err := p.Save(12*vg.Inch, 9*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

	getMeanSigma()
}

func getMeanSigma() {

	sum := 0.0
	for _, i := range deflected {
		sum += i.y
	}
	mean := sum / float64(len(deflected))
	diffSq := 0.0
	for _, i := range deflected {
		diffSq += math.Pow(i.y-mean, 2)
	}
	sigma := math.Sqrt(diffSq / float64(len(deflected)))
	fmt.Printf("Deflected Mean: %f\nDefelcted Sigma %f\n", mean, sigma)

	sum = 0.0
	for _, i := range undeflected {
		sum += i.y
	}
	mean = sum / float64(len(undeflected))
	diffSq = 0.0
	for _, i := range undeflected {
		diffSq += math.Pow(i.y-mean, 2)
	}
	sigma = math.Sqrt(diffSq / float64(len(undeflected)))
	fmt.Printf("Undeflected Mean: %f\nUndefelcted Sigma %f\n", mean, sigma)

}

// creates a plotter.XYs for the data with a time average of n data points
func getTimeAvgPoints(def bool, n int) plotter.XYs {

	if def {
		//initalize plotter.XYs with data_length/n points
		pts := make(plotter.XYs, len(deflected)/n)
		//index over all points
		for i := range pts {
			ysum := 0.0
			//sum up n data values for this points
			for j := 0; j < n; j++ {
				ysum += deflected[i*n+j].y
			}

			//add average to data point and center at middle time of averaged data points
			pts[i].X = deflected[i+(n/2)].x
			pts[i].Y = ysum / float64(n)
		}
		return pts
	}
	//see above now with undeflected
	pts := make(plotter.XYs, len(undeflected)/n)
	for i := range pts {
		ysum := 0.0
		for j := 0; j < n; j++ {
			ysum += undeflected[i*n+j].y
		}
		pts[i].X = undeflected[i+(n/2)].x
		pts[i].Y = ysum / float64(n)
	}
	return pts
}

// get gonumplot points from data
func getPoints(def bool) plotter.XYs {

	if def {
		pts := make(plotter.XYs, len(deflected))
		for i := range pts {
			pts[i].X = deflected[i].x
			pts[i].Y = deflected[i].y
		}
		return pts
	}
	pts := make(plotter.XYs, len(deflected))
	for i := range pts {
		pts[i].X = deflected[i].x
		pts[i].Y = deflected[i].y
	}
	return pts

}

// read files and make list
func readFile(cup bool, run int) {

	typeStr := ""

	if cup {
		typeStr = "cup"
	} else {
		typeStr = "faceplate"
	}

	//get all files in data directory
	files, err := os.ReadDir("data")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), typeStr+"_deflected") && strings.Contains(file.Name(), strconv.Itoa(run-1)) {
			f, err := os.Open("data/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)

			scanner.Scan()

			count := 0.0

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
				deflected = append(deflected, tmp)

			}

			break

		}
	}

	for _, file := range files {
		if strings.Contains(file.Name(), typeStr+"_undeflected") && strings.Contains(file.Name(), strconv.Itoa(run-1)) {
			f, err := os.Open("data/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)

			scanner.Scan()

			count := 0.0

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
				undeflected = append(undeflected, tmp)

			}

			break

		}
	}

}
