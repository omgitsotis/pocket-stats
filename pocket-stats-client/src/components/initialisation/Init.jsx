import React, {Component} from 'react';
import DatePicker from 'react-date-picker';

class InitSection extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: new Date();
        }
    }

    onChange = date => this.setState({ date });

    render() {
        return (
            <div>
                <DatePicker
                    onChange={this.onChange}
                    value={this.state.date}
            </div>
        );
    }
}
