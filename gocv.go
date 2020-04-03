package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

const MinimumArea = 3000

func main() {
	input, err := gocv.VideoCaptureFile("./movie/sample3.mp4")

	if err != nil {
		fmt.Printf("error opening video writer device: %v\n", err)
		return
	}
	defer input.Close()

	output, err := gocv.VideoWriterFile(
		"./movie/output3.mp4", 
		input.CodecString(), 
		input.Get(gocv.VideoCaptureFPS), 
		int(input.Get(gocv.VideoCaptureFrameWidth)), 
		int(input.Get(gocv.VideoCaptureFrameHeight)), 
		true)

	fmt.Println(input.Get(gocv.VideoCaptureFPS),input.Get(gocv.VideoCaptureFrameWidth),input.Get(gocv.VideoCaptureFrameHeight))

	if err != nil {
		fmt.Print(err)
		return
	}
	defer output.Close()

	img := gocv.NewMat()
	defer img.Close()
	imgDelta := gocv.NewMat()
	defer imgDelta.Close()
	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()

	//count := 0

	start := time.Now()

	for {
		if ok := input.Read(&img); !ok {
			fmt.Println("finish")
			break
		}
		if img.Empty() {
			continue
		}

		statusColor := color.RGBA{0, 255, 0, 0}
		mog2.Apply(img, &imgDelta)
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()

		gocv.Dilate(imgThresh, &imgThresh, kernel)

		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for i, c := range contours {
			area := gocv.ContourArea(c)

			if area < MinimumArea {
				continue
			}

			statusColor = color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(c)
			gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
		}

		output.Write(img)

		/*
				buf, _ := gocv.IMEncode(".jpg", img)

				filename := fmt.Sprintf("./images/sample%d.jpg", count)
				saveFile, err := os.Create(filename)

				if err != nil {
					fmt.Println(err)
				}
				defer saveFile.Close()

				saveFile.Write(buf)


			count++
		*/

	}
	end := time.Now()
	fmt.Println("Total time: ", end.Sub(start))
}
