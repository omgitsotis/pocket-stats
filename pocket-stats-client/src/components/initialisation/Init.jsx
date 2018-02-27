import React, {Component} from 'react';
import DatePicker from 'react-datepicker';
import moment from 'moment';

import 'react-datepicker/dist/react-datepicker.css';

class InitSection extends Component {
    constructor(props) {
        super(props);
        this.state = {
            date: moment()
        }
    }

    onChange = date => this.setState({ date });

    render() {
        const unix = this.state.date.startOf('day').unix();
        let section;

        switch (this.props.initState) {
            case 'started':
                section = <div>Loading, please wait</div>
                break;
            case 'completed':
                section =
                    <div>
                        <p>Completed</p>
                        <button onClick={this.props.onBackClick}>Return</button>
                    </div>
                    break;
            default:
                section =
                    <div>
                        <DatePicker
                            onChange={this.onChange}
                            selected={this.state.date}
                        />
                        <button onClick={() => this.props.onInitClick(unix)}>Initialise</button>
                        <button onClick={this.props.onBackClick}>Return</button>
                    </div>;
        }

        return (
            <div>
                {section}
            </div>
        );
    }
}

export default InitSection;
