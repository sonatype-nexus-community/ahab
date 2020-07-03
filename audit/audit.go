//
// Copyright 2019-Present Sonatype Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package audit

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/logrusorgru/aurora"
	"github.com/shopspring/decimal"
	"github.com/sonatype-nexus-community/go-sona-types/ossindex/types"
)

var (
	nine, seven, four decimal.Decimal
)

func init() {
	nine, _ = decimal.NewFromString("9")
	seven, _ = decimal.NewFromString("7")
	four, _ = decimal.NewFromString("4")
}

func LogResults(noColor bool, loud bool, output string, projects []types.Coordinate) (vulnerableCount int, results string, err error) {
	switch output {
	case "json":
		vulnerableCount, results, err = outputJSON(loud, projects)
	case "csv":
		vulnerableCount, results, err = outputCSV(loud, projects)
	default:
		vulnerableCount, results, err = outputText(noColor, loud, projects)
	}
	return
}

func outputJSON(loud bool, projects []types.Coordinate) (int, string, error) {
	_, vulnerablePackages := splitPackages(projects)
	if loud {
		b, err := json.Marshal(projects)
		if err != nil {
			fmt.Println(err)
		}
		return len(vulnerablePackages), string(b), nil
	}
	b, err := json.Marshal(vulnerablePackages)
	if err != nil {
		return 0, "", err
	}
	return len(vulnerablePackages), string(b), nil
}

func outputCSV(loud bool, projects []types.Coordinate) (int, string, error) {
	_, vulnerablePackages := splitPackages(projects)

	writer := func(w *csv.Writer, v types.Coordinate) (err error) {
		if v.IsVulnerable() {
			for _, vv := range v.Vulnerabilities {
				err = w.Write([]string{
					v.Coordinates,
					v.Reference,
					vv.ID,
					vv.Title,
					vv.Description,
					vv.Reference,
					vv.CvssScore.String(),
					vv.Cve,
					vv.CvssVector,
				})
				if err != nil {
					return
				}
			}
		} else {
			err = w.Write([]string{
				v.Coordinates,
				v.Reference,
			})
			if err != nil {
				return
			}
		}
		return
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	if loud {
		for _, v := range projects {
			err := writer(w, v)
			if err != nil {
				return 0, "", err
			}
		}
		w.Flush()
		return len(vulnerablePackages), b.String(), nil
	}

	for _, v := range vulnerablePackages {
		err := writer(w, v)
		if err != nil {
			return 0, "", err
		}
	}
	w.Flush()
	return len(vulnerablePackages), b.String(), nil
}

func outputText(noColor bool, loud bool, projects []types.Coordinate) (int, string, error) {
	var sb strings.Builder

	w := tabwriter.NewWriter(&sb, 9, 3, 0, '\t', 0)
	err := w.Flush()
	if err != nil {
		return 0, "", err
	}

	nonVulnerablePackages, vulnerablePackages := splitPackages(projects)

	err = groupAndPrint(vulnerablePackages, nonVulnerablePackages, loud, noColor, &sb)
	if err != nil {
		return 0, "", err
	}

	au := aurora.NewAurora(!noColor)
	t := table.NewWriter()
	t.SetStyle(table.StyleBold)
	t.SetTitle("Summary")
	t.AppendRow([]interface{}{"Audited Dependencies", strconv.Itoa(len(projects))})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"Vulnerable Dependencies", au.Bold(au.Red(strconv.Itoa(len(vulnerablePackages))))})
	sb.WriteString(t.Render())

	return len(vulnerablePackages), sb.String(), nil
}

func groupAndPrint(vulnerable []types.Coordinate, nonVulnerable []types.Coordinate, loud bool, noColor bool, sb *strings.Builder) (err error) {
	if loud {
		_, err = sb.WriteString("\nNon Vulnerable Packages\n\n")
		if err != nil {
			return
		}
		for k, v := range nonVulnerable {
			err = formatPackage(sb, noColor, k+1, len(nonVulnerable), v)
			if err != nil {
				return
			}
		}
	}
	if len(vulnerable) > 0 {
		_, err = sb.WriteString("\nVulnerable Packages\n\n")
		if err != nil {
			return
		}
		for k, v := range vulnerable {
			err = formatVulnerability(sb, noColor, k+1, len(vulnerable), v)
			if err != nil {
				return
			}
		}
	}

	return
}

func formatPackage(sb *strings.Builder, noColor bool, idx int, packageCount int, coordinate types.Coordinate) (err error) {
	au := aurora.NewAurora(!noColor)

	_, err = sb.WriteString(
		fmt.Sprintf("[%d/%d]\t%s\n",
			idx,
			packageCount,
			au.Bold(au.Green(coordinate.Coordinates)).String(),
		),
	)
	if err != nil {
		return
	}

	return
}

func formatVulnerability(sb *strings.Builder, noColor bool, idx int, packageCount int, coordinate types.Coordinate) (err error) {
	au := aurora.NewAurora(!noColor)

	_, err = sb.WriteString(fmt.Sprintf(
		"[%d/%d]\t%s\n%s \n",
		idx,
		packageCount,
		au.Bold(au.Red(coordinate.Coordinates)).String(),
		au.Red(strconv.Itoa(len(coordinate.Vulnerabilities))+" known vulnerabilities affecting installed version").String(),
	))
	if err != nil {
		return
	}
	sort.Slice(coordinate.Vulnerabilities, func(i, j int) bool {
		return coordinate.Vulnerabilities[i].CvssScore.GreaterThan(coordinate.Vulnerabilities[j].CvssScore)
	})

	for _, v := range coordinate.Vulnerabilities {
		if !v.Excluded {
			t := table.NewWriter()
			t.SetStyle(table.StyleBold)
			t.SetTitle(printColorBasedOnCvssScore(v.CvssScore, v.Title, noColor))
			t.AppendRow([]interface{}{"Description", text.WrapSoft(v.Description, 75)})
			t.AppendSeparator()
			t.AppendRow([]interface{}{"OSS Index ID", v.ID})
			t.AppendSeparator()
			t.AppendRow([]interface{}{"CVSS Score", fmt.Sprintf("%s/10 (%s)", v.CvssScore, scoreAssessment(v.CvssScore))})
			t.AppendSeparator()
			t.AppendRow([]interface{}{"CVSS Vector", v.CvssVector})
			t.AppendSeparator()
			t.AppendRow([]interface{}{"Link for more info", v.Reference})
			_, err = sb.WriteString(t.Render() + "\n")
			if err != nil {
				return
			}
		}
	}
	return
}

func printColorBasedOnCvssScore(score decimal.Decimal, text string, noColor bool) string {
	au := aurora.NewAurora(!noColor)
	if score.GreaterThanOrEqual(nine) {
		return au.Red(au.Bold(text)).String()
	}
	if score.GreaterThanOrEqual(seven) {
		return au.Red(text).String()
	}
	if score.GreaterThanOrEqual(four) {
		return au.Yellow(text).String()
	}
	return au.Green(text).String()
}

func scoreAssessment(score decimal.Decimal) string {
	if score.GreaterThanOrEqual(nine) {
		return "Critical"
	}
	if score.GreaterThanOrEqual(seven) {
		return "High"
	}
	if score.GreaterThanOrEqual(four) {
		return "Medium"
	}
	return "Low"
}

func splitPackages(entries []types.Coordinate) (nonVulnerable []types.Coordinate, vulnerable []types.Coordinate) {
	for _, v := range entries {
		if v.IsVulnerable() {
			vulnerable = append(vulnerable, v)
		} else {
			nonVulnerable = append(nonVulnerable, v)
		}
	}
	return
}
