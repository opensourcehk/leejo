package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
)

type RestTestResult interface {
	Count() int
	NthExists(int) error
	NthValid(int) error
	NthMatchPayload(int, *interface{}) error
}

type RestTester struct {
	BaseUrl string
}

func (t *RestTester) TestCreate(
	payload interface{}, result RestTestResult) (resp *napping.Response, err error) {

	// make the REST create request
	resp, err = napping.Post(t.BaseUrl,
		&payload, &result, nil)

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		err = fmt.Errorf("Bad response in TestCreate. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatchPayload(0, &payload)
	if err != nil {
		return
	}

	return
}

func (t *RestTester) TestRetrieveOne(
	id string, payload interface{}, result RestTestResult) (resp *napping.Response, err error) {

	// REST retrieve record with id
	p := napping.Params{} // empty payload for retrieve
	resp, err = napping.Get(t.BaseUrl+"/"+id,
		&p, &result, nil)
	if err != nil {
		return
	}

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		err = fmt.Errorf("Bad response in TestRetrieveOne. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatchPayload(0, &payload)
	if err != nil {
		return
	}

	return
}

func (t *RestTester) TestUpdate(
	id string, payload interface{}, result RestTestResult) (resp *napping.Response, err error) {

	// REST update record (of given id) with the payload
	resp, err = napping.Put(t.BaseUrl+"/"+id,
		&payload, &result, nil)
	if err != nil {
		return
	}

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		err = fmt.Errorf("Bad response in TestUpdate. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatchPayload(0, &payload)
	if err != nil {
		return
	}

	return
}

func (t *RestTester) TestDelete(
	id string, payload interface{}, result RestTestResult) (resp *napping.Response, err error) {

	// REST delete record of given id
	resp, err = napping.Delete(t.BaseUrl+"/"+id,
		&result, nil)
	if err != nil {
		return
	}

	// test: has to be exactly 1 result
	count := result.Count()
	if count != 1 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		err = fmt.Errorf("Bad response in TestUpdate. "+
			"There are %d results (expecting 1)",
			count)
		return
	}

	// test: test the result
	err = result.NthValid(0)
	if err != nil {
		return
	}
	err = result.NthMatchPayload(0, &payload)
	if err != nil {
		return
	}

	return
}
