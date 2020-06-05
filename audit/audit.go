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

func LogResults(quiet bool, noColor bool, loud bool, output string, projects []types.Coordinate) (vulnerableCount int, results string) {
	switch output {
	case "json":
		vulnerableCount, results = outputJSON(loud, projects)
	case "csv":
		vulnerableCount, results = outputCSV(loud, projects)
	default:
		vulnerableCount, results = outputText(quiet, noColor, loud, projects)
	}
	return
}

func outputJSON(loud bool, projects []types.Coordinate) (int, string) {
	_, vulnerablePackages := splitPackages(projects)
	if loud {
		b, err := json.Marshal(projects)
		if err != nil {
			fmt.Println(err)
		}
		return len(vulnerablePackages), string(b)
	}
	b, err := json.Marshal(vulnerablePackages)
	if err != nil {
		fmt.Println(err)
	}
	return len(vulnerablePackages), string(b)
}

func outputCSV(loud bool, projects []types.Coordinate) (int, string) {
	return 0, ""
}

func outputText(quiet bool, noColor bool, loud bool, projects []types.Coordinate) (int, string) {
	var sb strings.Builder

	w := tabwriter.NewWriter(&sb, 9, 3, 0, '\t', 0)
	w.Flush()
	var numVulnerable int

	nonVulnerablePackages, vulnerablePackages := splitPackages(projects)

	groupAndPrint(vulnerablePackages, nonVulnerablePackages, loud, noColor, &sb)
	au := aurora.NewAurora(!noColor)
	t := table.NewWriter()
	t.SetStyle(table.StyleBold)
	t.SetTitle("Summary")
	t.AppendRow([]interface{}{"Audited Dependencies", strconv.Itoa(len(projects))})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"Vulnerable Dependencies", au.Bold(au.Red(strconv.Itoa(numVulnerable)))})
	sb.WriteString(t.Render())

	return len(vulnerablePackages), sb.String()
}

func groupAndPrint(vulnerable []types.Coordinate, nonVulnerable []types.Coordinate, loud bool, noColor bool, sb *strings.Builder) {
	if loud {
		sb.WriteString("\nNon Vulnerable Packages\n\n")
		for k, v := range nonVulnerable {
			formatPackage(sb, noColor, k+1, len(nonVulnerable), v)
		}
	}
	if len(vulnerable) > 0 {
		sb.WriteString("\nVulnerable Packages\n\n")
		for k, v := range vulnerable {
			formatVulnerability(sb, noColor, k+1, len(vulnerable), v)
		}
	}
}

func formatPackage(sb *strings.Builder, noColor bool, idx int, packageCount int, coordinate types.Coordinate) {
	au := aurora.NewAurora(!noColor)

	sb.WriteString(
		fmt.Sprintf("[%d/%d]\t%s\n",
			idx,
			packageCount,
			au.Bold(au.Green(coordinate.Coordinates)).String(),
		),
	)
}

func formatVulnerability(sb *strings.Builder, noColor bool, idx int, packageCount int, coordinate types.Coordinate) {
	au := aurora.NewAurora(!noColor)

	sb.WriteString(fmt.Sprintf(
		"[%d/%d]\t%s\n%s \n",
		idx,
		packageCount,
		au.Bold(au.Red(coordinate.Coordinates)).String(),
		au.Red(strconv.Itoa(len(coordinate.Vulnerabilities))+" known vulnerabilities affecting installed version").String(),
	))
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
			sb.WriteString(t.Render() + "\n")
		}
	}
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
