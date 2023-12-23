package main

import (
	"context"
	"log"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	SubscriptionStatistics struct {
		Timeframe           string
		ActiveSubscription  int64
		BytesSubscription   int64
		EndSubscription     int64
		HoursSubscription   int64
		StartSubscription   int64
		SubscriptionBytes   string
		SubscriptionDeposit types.Coins
		SubscriptionHours   int64
		SubscriptionRefund  types.Coins
	}
)

func NewSubscriptionStatistics(timeframe string) *SubscriptionStatistics {
	return &SubscriptionStatistics{
		Timeframe:           timeframe,
		SubscriptionDeposit: types.NewCoins(nil),
		SubscriptionRefund:  types.NewCoins(nil),
	}
}

func (s *SubscriptionStatistics) Result(addr string, timestamp time.Time) bson.M {
	res := bson.M{
		"addr":      addr,
		"timeframe": s.Timeframe,
		"timestamp": timestamp,
	}

	if s.ActiveSubscription != 0 {
		res["active_subscription"] = s.ActiveSubscription
	}
	if s.BytesSubscription != 0 {
		res["bytes_subscription"] = s.BytesSubscription
	}
	if s.EndSubscription != 0 {
		res["end_subscription"] = s.EndSubscription
	}
	if s.HoursSubscription != 0 {
		res["hours_subscription"] = s.HoursSubscription
	}
	if s.StartSubscription != 0 {
		res["start_subscription"] = s.StartSubscription
	}
	if s.SubscriptionBytes != "" && s.SubscriptionBytes != "0" {
		res["subscription_bytes"] = s.SubscriptionBytes
	}
	if s.SubscriptionDeposit.Len() != 0 {
		res["subscription_deposit"] = s.SubscriptionDeposit
	}
	if s.SubscriptionHours != 0 {
		res["subscription_hours"] = s.SubscriptionHours
	}
	if s.SubscriptionRefund.Len() != 0 {
		res["subscription_refund"] = s.SubscriptionRefund
	}

	return res
}

func StatisticsFromSubscriptions(ctx context.Context, db *mongo.Database, minTimestamp, maxTimestamp time.Time, excludeAddrs []string) (result []bson.M, err error) {
	log.Println("StatisticsFromSubscriptions", minTimestamp, maxTimestamp)

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"acc_addr":        1,
		"end_timestamp":   1,
		"deposit":         1,
		"gigabytes":       1,
		"hours":           1,
		"node_addr":       1,
		"refund":          1,
		"start_timestamp": 1,
	}

	items, err := database.SubscriptionFind(ctx, db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[string]map[time.Time]*SubscriptionStatistics)
		w = make(map[string]map[time.Time]*SubscriptionStatistics)
		m = make(map[string]map[time.Time]*SubscriptionStatistics)
		y = make(map[string]map[time.Time]*SubscriptionStatistics)
	)

	for i := 0; i < len(items); i++ {
		if _, ok := d[items[i].NodeAddr]; !ok {
			d[items[i].NodeAddr] = make(map[time.Time]*SubscriptionStatistics)
		}
		if _, ok := w[items[i].NodeAddr]; !ok {
			w[items[i].NodeAddr] = make(map[time.Time]*SubscriptionStatistics)
		}
		if _, ok := m[items[i].NodeAddr]; !ok {
			m[items[i].NodeAddr] = make(map[time.Time]*SubscriptionStatistics)
		}
		if _, ok := y[items[i].NodeAddr]; !ok {
			y[items[i].NodeAddr] = make(map[time.Time]*SubscriptionStatistics)
		}

		exclude := utils.ContainsString(excludeAddrs, items[i].AccAddr)
		if exclude {
			continue
		}

		startTimestamp := items[i].StartTimestamp
		if items[i].StartTimestamp.IsZero() {
			startTimestamp = minTimestamp
		}

		endTimestamp := items[i].EndTimestamp
		if items[i].EndTimestamp.IsZero() {
			endTimestamp = maxTimestamp
		}

		dayStartTimestamp, dayEndTimestamp := utils.DayDate(startTimestamp), utils.DayDate(endTimestamp)
		for t := dayStartTimestamp; !t.After(dayEndTimestamp); t = t.AddDate(0, 0, 1) {
			if _, ok := d[items[i].NodeAddr][t]; !ok {
				d[items[i].NodeAddr][t] = NewSubscriptionStatistics("day")
			}

			d[items[i].NodeAddr][t].ActiveSubscription += 1
		}

		weekStartTimestamp, weekEndTimestamp := utils.ISOWeekDate(startTimestamp), utils.ISOWeekDate(endTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			if _, ok := w[items[i].NodeAddr][t]; !ok {
				w[items[i].NodeAddr][t] = NewSubscriptionStatistics("week")
			}

			w[items[i].NodeAddr][t].ActiveSubscription += 1
		}

		monthStartTimestamp, monthEndTimestamp := utils.MonthDate(startTimestamp), utils.MonthDate(endTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			if _, ok := m[items[i].NodeAddr][t]; !ok {
				m[items[i].NodeAddr][t] = NewSubscriptionStatistics("month")
			}

			m[items[i].NodeAddr][t].ActiveSubscription += 1
		}

		yearStartTimestamp, yearEndTimestamp := utils.YearDate(startTimestamp), utils.YearDate(endTimestamp)
		for t := yearStartTimestamp; !t.After(yearEndTimestamp); t = t.AddDate(1, 0, 0) {
			if _, ok := y[items[i].NodeAddr][t]; !ok {
				y[items[i].NodeAddr][t] = NewSubscriptionStatistics("year")
			}

			y[items[i].NodeAddr][t].ActiveSubscription += 1
		}

		if !items[i].EndTimestamp.IsZero() {
			d[items[i].NodeAddr][dayEndTimestamp].EndSubscription += 1
			w[items[i].NodeAddr][weekEndTimestamp].EndSubscription += 1
			m[items[i].NodeAddr][monthEndTimestamp].EndSubscription += 1
			y[items[i].NodeAddr][yearEndTimestamp].EndSubscription += 1
		}
		if items[i].Gigabytes != 0 {
			d[items[i].NodeAddr][dayStartTimestamp].BytesSubscription += 1
			w[items[i].NodeAddr][weekStartTimestamp].BytesSubscription += 1
			m[items[i].NodeAddr][monthStartTimestamp].BytesSubscription += 1
			y[items[i].NodeAddr][yearStartTimestamp].BytesSubscription += 1

			bytes := hubtypes.Gigabyte.MulRaw(items[i].Gigabytes)
			d[items[i].NodeAddr][dayStartTimestamp].SubscriptionBytes = utils.MustIntFromString(d[items[i].NodeAddr][dayStartTimestamp].SubscriptionBytes).Add(bytes).String()
			w[items[i].NodeAddr][weekStartTimestamp].SubscriptionBytes = utils.MustIntFromString(w[items[i].NodeAddr][weekStartTimestamp].SubscriptionBytes).Add(bytes).String()
			m[items[i].NodeAddr][monthStartTimestamp].SubscriptionBytes = utils.MustIntFromString(m[items[i].NodeAddr][monthStartTimestamp].SubscriptionBytes).Add(bytes).String()
			y[items[i].NodeAddr][yearStartTimestamp].SubscriptionBytes = utils.MustIntFromString(y[items[i].NodeAddr][yearStartTimestamp].SubscriptionBytes).Add(bytes).String()
		}
		if items[i].Hours != 0 {
			d[items[i].NodeAddr][dayStartTimestamp].HoursSubscription += 1
			w[items[i].NodeAddr][weekStartTimestamp].HoursSubscription += 1
			m[items[i].NodeAddr][monthStartTimestamp].HoursSubscription += 1
			y[items[i].NodeAddr][yearStartTimestamp].HoursSubscription += 1

			d[items[i].NodeAddr][dayStartTimestamp].SubscriptionHours += items[i].Hours
			w[items[i].NodeAddr][weekStartTimestamp].SubscriptionHours += items[i].Hours
			m[items[i].NodeAddr][monthStartTimestamp].SubscriptionHours += items[i].Hours
			y[items[i].NodeAddr][yearStartTimestamp].SubscriptionHours += items[i].Hours
		}
		if !items[i].StartTimestamp.IsZero() {
			d[items[i].NodeAddr][dayStartTimestamp].StartSubscription += 1
			w[items[i].NodeAddr][weekStartTimestamp].StartSubscription += 1
			m[items[i].NodeAddr][monthStartTimestamp].StartSubscription += 1
			y[items[i].NodeAddr][yearStartTimestamp].StartSubscription += 1
		}
	}

	for i := 0; i < len(items); i++ {
		startTimestamp := items[i].StartTimestamp
		if items[i].StartTimestamp.IsZero() {
			startTimestamp = minTimestamp
		}

		endTimestamp := items[i].EndTimestamp
		if items[i].EndTimestamp.IsZero() {
			endTimestamp = maxTimestamp
		}

		dayStartTimestamp, dayEndTimestamp := utils.DayDate(startTimestamp), utils.DayDate(endTimestamp)
		if _, ok := d[items[i].NodeAddr][dayStartTimestamp]; !ok {
			d[items[i].NodeAddr][dayStartTimestamp] = NewSubscriptionStatistics("day")
		}
		if _, ok := d[items[i].NodeAddr][dayEndTimestamp]; !ok {
			d[items[i].NodeAddr][dayEndTimestamp] = NewSubscriptionStatistics("day")
		}

		weekStartTimestamp, weekEndTimestamp := utils.ISOWeekDate(startTimestamp), utils.ISOWeekDate(endTimestamp)
		if _, ok := w[items[i].NodeAddr][weekStartTimestamp]; !ok {
			w[items[i].NodeAddr][weekStartTimestamp] = NewSubscriptionStatistics("week")
		}
		if _, ok := w[items[i].NodeAddr][weekEndTimestamp]; !ok {
			w[items[i].NodeAddr][weekEndTimestamp] = NewSubscriptionStatistics("week")
		}

		monthStartTimestamp, monthEndTimestamp := utils.MonthDate(startTimestamp), utils.MonthDate(endTimestamp)
		if _, ok := m[items[i].NodeAddr][monthStartTimestamp]; !ok {
			m[items[i].NodeAddr][monthStartTimestamp] = NewSubscriptionStatistics("month")
		}
		if _, ok := m[items[i].NodeAddr][monthEndTimestamp]; !ok {
			m[items[i].NodeAddr][monthEndTimestamp] = NewSubscriptionStatistics("month")
		}

		yearStartTimestamp, yearEndTimestamp := utils.YearDate(startTimestamp), utils.YearDate(endTimestamp)
		if _, ok := y[items[i].NodeAddr][yearStartTimestamp]; !ok {
			y[items[i].NodeAddr][yearStartTimestamp] = NewSubscriptionStatistics("year")
		}
		if _, ok := y[items[i].NodeAddr][yearEndTimestamp]; !ok {
			y[items[i].NodeAddr][yearEndTimestamp] = NewSubscriptionStatistics("year")
		}

		if items[i].Deposit != nil {
			d[items[i].NodeAddr][dayStartTimestamp].SubscriptionDeposit = d[items[i].NodeAddr][dayStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			w[items[i].NodeAddr][weekStartTimestamp].SubscriptionDeposit = w[items[i].NodeAddr][weekStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			m[items[i].NodeAddr][monthStartTimestamp].SubscriptionDeposit = m[items[i].NodeAddr][monthStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			y[items[i].NodeAddr][yearStartTimestamp].SubscriptionDeposit = y[items[i].NodeAddr][yearStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
		}
		if items[i].Refund != nil {
			d[items[i].NodeAddr][dayEndTimestamp].SubscriptionRefund = d[items[i].NodeAddr][dayEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
			w[items[i].NodeAddr][weekEndTimestamp].SubscriptionRefund = w[items[i].NodeAddr][weekEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
			m[items[i].NodeAddr][monthEndTimestamp].SubscriptionRefund = m[items[i].NodeAddr][monthEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
			y[items[i].NodeAddr][yearEndTimestamp].SubscriptionRefund = y[items[i].NodeAddr][yearEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
		}
	}

	for s := range d {
		for t := range d[s] {
			result = append(result, d[s][t].Result(s, t))
		}
	}

	for s := range w {
		for t := range w[s] {
			result = append(result, w[s][t].Result(s, t))
		}
	}

	for s := range m {
		for t := range m[s] {
			result = append(result, m[s][t].Result(s, t))
		}
	}

	for s := range y {
		for t := range y[s] {
			result = append(result, y[s][t].Result(s, t))
		}
	}

	return result, nil
}
