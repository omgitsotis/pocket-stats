import React, {Component} from 'react';
import BoxScore from './BoxScore.jsx';
import Timeframe from './Timeframe.jsx';
import moment from 'moment';

class BoxScoreContainer extends Component {
    constructor(props) {
        super(props);
        this.state = {
          interval: "week",
          startDate: moment().utc().startOf("day").subtract(1, "w").unix(),
          endDate: moment().utc().startOf("day").unix(),
        };
    }

    onChange(e) {
      const interval = e.target.value
      this.setState({
        interval: interval
      });

      let startDate, endDate;

      switch (interval) {
        case "week":
          endDate = moment().utc().startOf("day");
          startDate = moment().utc().startOf("day").subtract(1, "M");
          break;

        case "month":
          endDate = moment().utc().startOf("day");
          startDate = moment().utc().startOf("day").subtract(1, "M");
          break;

        case "thirty":
          endDate = moment().utc().startOf("day");
          startDate = moment().utc().startOf("day").subtract(30, "d");
          break;

        case "sixty":
          endDate = moment().utc().startOf("day");
          startDate = moment().utc().startOf("day").subtract(60, "d");
          break;

        case "year":
          endDate = moment().utc().startOf("day");
          startDate = moment().utc().startOf("day").subtract(1, "y");
          break;

        default:
          return;
      }

      this.props.onFetchDataClick(startDate.unix(), endDate.unix())
    }

    onDateChange(date, value) {
      this.setState({
        [date]: value.startOf("day").unix()
      });
    }

    onCustomRequest() {
      this.props.onFetchDataClick(this.state.startDate, this.state.endDate)
    }

    render() {
      return (
        <div>
          <Timeframe
            interval={this.state.interval}
            startDate={this.state.startDate}
            endDate={this.state.endDate}
            onChange={(e) => this.onChange(e)}
            onDateChange={(date, value) => this.onDateChange(date, value)}
            onCustomRequest={() => this.onCustomRequest()} />
          <BoxScore totals={this.props.totals}/>
        </div>
      )
    }
}

export default BoxScoreContainer;
