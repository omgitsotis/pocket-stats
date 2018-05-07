import React, {Component} from 'react';

import Homepage from '../homepage/Homepage.jsx';
import Navbar from '../navbar/Navbar.jsx';
import BoxScoreContainer from '../boxscore/BoxScoreContainer.jsx';

class DashboardContainer extends Component {
    constructor(props) {
        super(props);
        this.state = {
            currentPage: 'home'
        };
    }

    onNavbarClick(state) {
        this.setState({currentPage: state});
    }

    render() {
        let component;
        switch (this.state.currentPage) {
            case 'home':
                component = <Homepage {...this.props} />;
                break;
            case 'boxscore':
              component =
                <BoxScoreContainer
                  totals={this.props.totals}
                  onFetchDataClick={this.props.onFetchDataClick} />;
              break;
            default:
                break;
        }

        // let component = this.props.isDataError ?
        //     <div>There was an error</div> :

        return (
            <div>
                <Navbar
                    updateComplete={this.props.updateComplete}
                    lastUpdated={this.props.lastUpdated}
                    onUpdateClick={this.props.onUpdateClick}
                    currentPage={this.state.currentPage}
                    onNavbarClick={(state) => this.onNavbarClick(state)}
                />
                {component}
            </div>
        );
    }
}

export default DashboardContainer;
