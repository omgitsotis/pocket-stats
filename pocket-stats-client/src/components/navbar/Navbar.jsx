import React, {Component} from 'react';
import moment from 'moment';
import classnames from 'classnames';

class Navbar extends Component {
    render() {
        const lastUpdated = moment.unix(this.props.lastUpdated).format("D/MMM");
        const currentDate = moment().startOf('day').unix();

        const isDisabled = (!this.props.updateComplete ||
            currentDate === this.props.lastUpdated);

        const btnClass = classnames({
            'btn': true,
            'btn-primary': !isDisabled,
            'btn-disabled': isDisabled
        });

        const iconClass = classnames({
            'fa': true,
            'fa-refresh': true,
            'fa-spin': !this.props.updateComplete
        });

        const state = 'home';

        return (
            <div className='row'>
                <div className='col-lg'>
                    <nav className="navbar navbar-expand-lg navbar-light bg-light">
                        <span className="navbar-brand mb-0 h1">Navbar</span>
                        <ul className="navbar-nav mr-auto">
                            <li className="nav-item">
                                <button
                                    type="button"
                                    class="btn btn-light nav-link"
                                    disabled={state === 'home'}>
                                    Home
                                </button>
                            </li>
                            <li className="nav-item active">
                                <button
                                    type="button"
                                    class="btn btn-light nav-link"
                                    disabled={state === 'box-score'}>
                                    Box Score
                                </button>
                            </li>
                            <li className="nav-item active">
                                <button
                                    type="button"
                                    class="btn btn-light nav-link"
                                    disabled={state == 'graph'}>
                                    Graphs
                                </button>
                            </li>
                        </ul>
                        <div>
                            <span className="navbar-text">Updated: {lastUpdated}</span>
                            <button type="button"
                                id="update-btn"
                                className={btnClass}
                                disabled={isDisabled}
                                onClick={() => this.props.onUpdateClick()}>
                                <i className={iconClass} aria-hidden="true"></i>
                            </button>
                        </div>
                    </nav>
                </div>
            </div>
        )
    }
}

export default Navbar;