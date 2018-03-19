import React, {Component} from 'react';
import {Line} from 'react-chartjs-2';
import moment from 'moment';
import classnames from 'classnames';
import './dashboard.css';


class Dashboard extends Component {
    render() {
        const totals = this.props.totals;
        const itemised = this.props.itemised;

        const lastUpdated = moment.unix(this.props.lastUpdated).format("D/MMM");
        const currentDate = moment().startOf('day').unix();

        const isDisabled = (!this.props.updateComplete || currentDate === this.props.lastUpdated)
        console.log(currentDate, "===", this.props.lastUpdated);
        console.log(!this.props.updateComplete);
        console.log(isDisabled);

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

        let labels = [];
        let atsRead = [];
        let atsAdded = [];
        Object.keys(itemised).forEach(function(key) {
            let day = moment.unix(key);
            labels.push(day.format("D/MMM"));
            atsRead.push(itemised[key].articles_read);
            atsAdded.push(itemised[key].articles_added);
        });

        const data = {
            labels: labels,
            datasets: [{
                label: 'Articles Read',
                fill: false,
                lineTension: 0.1,
                backgroundColor: 'rgba(75,192,192,0.4)',
                borderColor: 'rgba(75,192,192,1)',
                borderCapStyle: 'butt',
                borderDash: [],
                borderDashOffset: 0.0,
                borderJoinStyle: 'miter',
                pointBorderColor: 'rgba(75,192,192,1)',
                pointBackgroundColor: '#fff',
                pointBorderWidth: 1,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: 'rgba(75,192,192,1)',
                pointHoverBorderColor: 'rgba(220,220,220,1)',
                pointHoverBorderWidth: 2,
                pointRadius: 1,
                pointHitRadius: 10,
                data: atsRead
            },
            {
                label: 'Articles Added',
                fill: false,
                lineTension: 0.1,
                backgroundColor: 'rgba(128,0,0,0.4)',
                borderColor: 'rgba(128,0,0,1)',
                borderCapStyle: 'butt',
                borderDash: [],
                borderDashOffset: 0.0,
                borderJoinStyle: 'miter',
                pointBorderColor: 'rgba(75,192,192,1)',
                pointBackgroundColor: '#fff',
                pointBorderWidth: 1,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: 'rgba(128,0,0,1)',
                pointHoverBorderColor: 'rgba(128,0,0,1)',
                pointHoverBorderWidth: 2,
                pointRadius: 1,
                pointHitRadius: 10,
                data: atsAdded
            }]
        };

        return(
            <div>
                <div className='row'>
                    <div className='col-lg'>
                        <nav className="navbar navbar-light bg-light">
                            <span className="navbar-brand mb-0 h1">Navbar</span>
                            <div>
                                <span className="navbar-text">Updated: {lastUpdated}</span>
                                <button type="button"
                                    className={btnClass}
                                    disabled={isDisabled}
                                    onClick={() => this.props.onUpdateClick()}>
                                    <i className={iconClass} aria-hidden="true"></i>
                                </button>
                            </div>
                        </nav>
                    </div>
                </div>
                <div className='row'>
                    <div className="col-lg">
                        <div className='card-deck'>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Articles Add</h5>
                                <p className='card-text'>{totals.total_articles_added}</p>
                            </div>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Articles Read</h5>
                                <p className='card-text'>{totals.total_articles_read}</p>
                            </div>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Words add</h5>
                                <p className='card-text'>{totals.total_words_added}</p>
                            </div>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Words Read</h5>
                                <p className='card-text'>{totals.total_words_read}</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="row">
                    <div className='col-lg'>
                        <Line data={data} height={100}/>
                    </div>
                </div>
            </div>
        );
    }
}

export default Dashboard;
