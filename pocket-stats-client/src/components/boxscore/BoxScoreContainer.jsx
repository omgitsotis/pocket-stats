import React, {Component} from 'react';
import BoxScore from './BoxScore.jsx';
import Timeframe from './Timeframe.jsx';
import moment from 'moment';

class BoxScoreContainer extends Component {
    constructor(props) {
        super(props);
        this.state = {
          interval: "week",
          startDate: moment().utc().startOf("day").subtract(1, "month").unix(),
          endDate: moment().utc().startOf("day").unix(),
        };
    }

    onChange(e) {
      this.setState({
        interval: e.target.value
      });
    }

    onDateChange(date, value) {
      console.log(date, value)
      this.setState({
        [date]: value.startOf("day").unix()
      });
    }

    render() {
      return (
        <div>
          <Timeframe
            onChange={ (e) => this.onChange(e) }
            interval={this.state.interval}
            onDateChange={ (date, value) => this.onDateChange(date, value) }
            startDate={this.state.startDate}/>
          <BoxScore totals={this.props.totals}/>
        </div>
      )
    }
}

export default BoxScoreContainer;
