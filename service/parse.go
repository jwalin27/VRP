package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Vorto-interview/models"
	"github.com/pkg/errors"
)

// pointMatcher is a regex to match parentheses-surrounded x,y coordinates of floating point numbers. Ex: `(1.5,2.9)`
var pointMatcher = regexp.MustCompile(`\((-?\d+\.?\d*),(-?\d+\.?\d*)\)`)

// ParseAllLoads reads the file at the provided path and parses out a slice of Loads
func ParseAllLoads(inFile string) ([]models.Load, error) {
	file, err := os.Open(inFile)
	if err != nil {
		return []models.Load{}, errors.Wrap(err, "failed to read load file")
	}
	defer file.Close()

	return parseReader(file)
}

func parseReader(reader io.Reader) ([]models.Load, error) {
	loads := make([]models.Load, 0)

	scanner := bufio.NewScanner(reader)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum == 1 {
			// Skip the header line
			continue
		}

		l, err := parseLine(scanner.Text())
		if err != nil {
			return []models.Load{}, errors.Wrap(err, fmt.Sprintf("unable to parse line %v", lineNum))
		}

		loads = append(loads, l)
	}

	return loads, nil
}

func parseLine(input string) (models.Load, error) {
	fields := strings.Fields(input)
	if len(fields) != 3 {
		return models.Load{}, fmt.Errorf("malformed load: expected 3 fields, received %v", len(fields))
	}

	loadNum, err := strconv.Atoi(fields[0])
	if err != nil {
		return models.Load{}, errors.Wrap(err, "load number unable to be parsed")
	}

	pickup, err := parsePoint(fields[1])
	if err != nil {
		return models.Load{}, errors.Wrap(err, "pickup point unable to be parsed")
	}

	dropoff, err := parsePoint(fields[2])
	if err != nil {
		return models.Load{}, errors.Wrap(err, "dropoff point unable to be parsed")
	}

	return models.Load{
		Number:  loadNum,
		Pickup:  pickup,
		Dropoff: dropoff,
	}, nil
}

func parsePoint(input string) (models.Point, error) {
	if !pointMatcher.MatchString(input) {
		return models.Point{}, fmt.Errorf("point does not match the regex %v", pointMatcher.String())
	}

	groups := pointMatcher.FindStringSubmatch(input)

	x, err := strconv.ParseFloat(groups[1], 64)
	if err != nil {
		return models.Point{}, errors.Wrap(err, "point x coordinate unable to be parsed")
	}

	y, err := strconv.ParseFloat(groups[2], 64)
	if err != nil {
		return models.Point{}, errors.Wrap(err, "point y coordinate unable to be parsed")
	}

	return models.Point{
		X: x,
		Y: y,
	}, nil
}
