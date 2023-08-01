package prompts

import (
	"github.com/eurozulu/pempal/ui"
	"strconv"
	"strings"
	"time"
)

const defaultFormat = time.RFC850

var monthNames = []string{
	"",
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

const yearDuration = 365 * (time.Hour * 24)

type DateView struct {
	ui.TextList
	Format string
}

func (dv DateView) Render(frame ui.ViewFrame) {
	dv.setChildValues(dv.String())
	dv.TextList.Render(frame)
}

func (dv *DateView) OnChildUpdate(child ui.View) {
	dv.SetText(dv.getChildValues())
}

func (dv *DateView) setChildValues(date string) {
	tm, err := time.Parse(dv.format(), date)
	if err != nil {
		dv.SetText(err.Error())
		tm = time.Now()
	}
	for _, cv := range dv.ChildViews() {
		if tv, ok := cv.(ui.TextView); ok {
			dv.setChildValue(tv, tm)
		}
	}
}

func (dv *DateView) setChildValue(child ui.TextView, t time.Time) {
	switch child.Label() {
	case "Year":
		child.SetText(strconv.Itoa(t.Year()))
	case "Month":
		child.SetText(t.Month().String())
	case "Day":
		child.SetText(strconv.Itoa(t.Day()))
	case "Hour":
		child.SetText(strconv.Itoa(t.Hour()))
	case "Minute":
		child.SetText(strconv.Itoa(t.Minute()))
	case "Second":
		child.SetText(strconv.Itoa(t.Second()))
	default:
		// ignore it
	}
}

func (dv *DateView) getChildValues() string {
	year := dv.getChildValueAsInt("Year")
	month := dv.getChildValueAsMonth("Month")
	day := dv.getChildValueAsInt("Day")
	hour := dv.getChildValueAsInt("Hour")
	minute := dv.getChildValueAsInt("Minute")
	second := dv.getChildValueAsInt("Second")
	tm := time.Date(year, month, day, hour, minute, second, 0, time.Local)
	return tm.Format(dv.format())
}

func (dv DateView) getChildValueAsMonth(label string) time.Month {
	v := dv.ChildByLabel(label)
	if v == nil || v.String() == "" {
		return 0
	}
	return parseMonth(v.String())
}

func (dv DateView) getChildValueAsInt(label string) int {
	v := dv.ChildByLabel(label)
	if v == nil {
		return 0
	}
	i, err := strconv.Atoi(v.String())
	if err != nil {
		return 0
	}
	return i
}

func (dv DateView) format() string {
	if dv.Format == "" {
		return defaultFormat
	}
	return dv.Format
}

func parseMonth(s string) time.Month {
	for i, mn := range monthNames {
		if strings.EqualFold(s, mn) {
			return time.Month(i)
		}
	}
	return 0
}

func rangeList(begin, length int) []string {
	var rangeValues []string
	for i := begin; i < begin+length; i++ {
		rangeValues = append(rangeValues, strconv.Itoa(i))
	}
	return rangeValues
}

func buildDateChildViews() []ui.View {
	now := time.Now()
	return []ui.View{
		ui.NewTextList("Year", "", rangeList(now.Year(), now.Add(yearDuration*25).Year())...),
		ui.NewTextList("Month", "", monthNames[1:]...),
		ui.NewTextList("Day", "", rangeList(1, 31)...),
		ui.NewTextList("Hour", "", rangeList(0, 23)...),
		ui.NewTextList("Minute", "", rangeList(0, 59)...),
		ui.NewTextList("Second", "", rangeList(0, 59)...),
	}
}

func NewDateView(label, date string) *DateView {
	return &DateView{TextList: *ui.NewTextListView(label, date, buildDateChildViews()...)}
}
