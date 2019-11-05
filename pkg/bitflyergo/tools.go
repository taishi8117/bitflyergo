package bitflyergo

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func (bf *Bitflyer) GetChildOrdersByDate(productCode string, from time.Time, to time.Time) ([]ChildOrder, error) {

	params := map[string]string{}
	params["product_code"] = productCode
	params["count"] = "500"
	var after int64 = 0
	var tmp []ChildOrder

	for {

		// set params
		if after > 0 {
			params["before"] = strconv.Itoa(int(after))
		}

		// fetch
		executions, err := bf.GetChildOrders(params)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(executions) == 0 {
			break
		}

		// sort by id asc
		sort.Slice(executions, func(i, j int) bool {
			return executions[i].Id < executions[j].Id
		})

		tmp = append(tmp, executions...)

		if executions[0].ChildOrderDate.Before(from) {
			break
		} else {
			after = executions[0].Id
		}
		time.Sleep(1 * time.Second)
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Id < tmp[j].Id
	})

	start := 0
	end := len(tmp)
	for i, e := range tmp {

		// find start index
		if start == 0 && (e.ChildOrderDate.Equal(from) || e.ChildOrderDate.After(from)) {
			//fmt.Println("start", i, e.ExecDate)
			start = i
			continue
		}

		// find end index
		if e.ChildOrderDate.After(to) {
			//fmt.Println("end", i, e.ExecDate)
			end = i
			break
		}
	}

	return tmp[start:end], nil
}

func (bf *Bitflyer) GetMyExecutionsByDate(productCode string, from time.Time, to time.Time) ([]MyExecution, error) {

	params := map[string]string{}
	params["product_code"] = productCode
	params["count"] = "500"
	var after int64 = 0
	var tmp []MyExecution

	for {

		// set params
		if after > 0 {
			params["before"] = strconv.Itoa(int(after))
		}

		// fetch
		executions, err := bf.GetMyExecutions(params)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(executions) == 0 {
			break
		}

		// sort by id asc
		sort.Slice(executions, func(i, j int) bool {
			return executions[i].Id < executions[j].Id
		})

		tmp = append(tmp, executions...)

		if executions[0].ExecDate.Before(from) {
			break
		} else {
			after = executions[0].Id
		}
		time.Sleep(1 * time.Second)
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Id < tmp[j].Id
	})

	start := 0
	end := len(tmp)
	for i, e := range tmp {

		// find start index
		if start == 0 && (e.ExecDate.Equal(from) || e.ExecDate.After(from)) {
			//fmt.Println("start", i, e.ExecDate)
			start = i
			continue
		}

		// find end index
		if e.ExecDate.After(to) {
			//fmt.Println("end", i, e.ExecDate)
			end = i
			break
		}
	}

	return tmp[start:end], nil
}

func (bf *Bitflyer) GetRelatedExecutionByOrder(order ChildOrder) ([]MyExecution, error) {
	params := map[string]string{}
	params["product_code"] = order.ProductCode
	params["child_order_acceptance_id"] = order.ChildOrderAcceptanceId
	executions, err := bf.GetMyExecutions(params)
	if err != nil {
		return nil, err
	}
	return executions, nil
}
