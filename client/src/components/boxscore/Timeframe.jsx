import React, {Component} from 'react'
import DatePicker from 'react-datepicker';
import moment from 'moment';


function Timeframe(props) {
  return(
    <div className="row timeframe">
      <div className="col-lg">
        <select
          className="interval"
          onChange={props.onChange}
          defaultValue={props.interval}>
          <option value="week">This week</option>
          <option value="month">This month</option>
          <option value="thirty">Last 30 days</option>
          <option value="sixty">Last 60 days</option>
          <option value="year">This year</option>
          <option value="custom"> Custom date</option>
        </select>
      </div>
      {props.interval === "custom" && <DatePickers {...props} />}
    </div>
  )
}

function DatePickers(props) {
  return(
    <div className="col-lg">
      <div className="row">
        <div className="col-md-4 offset-md-1">
          <DatePicker
            className="date-picker"
            selected={moment.unix(props.startDate)}
            onChange={(date) => props.onDateChange("startDate", date)} />
        </div>
        <div className="col-md-4 date-picker">
          <DatePicker
            selected={moment.unix(props.endDate)}
            onChange={(date) => props.onDateChange("endDate", date)} />
        </div>
        <div className="col-md-3">
        <button
            type="button"
            className="btn btn-light get-bttn"
            onClick={props.onCustomRequest}>
            Update
        </button>
        </div>
      </div>
    </div>
  )
}

export default Timeframe
