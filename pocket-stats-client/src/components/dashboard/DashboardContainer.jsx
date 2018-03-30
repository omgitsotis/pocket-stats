import React, {Component} from 'react';
import Dashboard from './Dashboard.jsx';

class DashboardContainer extends Component {
    constructor(props) {
        super(props);
        this.state = {
            currentPage: 'home'
        };
    }

    render() {
        let component = this.props.isDataError ?
            <div>There was an error</div> :
            <Dashboard {...this.props} />;

        return (
            <div>
                {component}
            </div>
        );
    }
}

export default DashboardContainer;
