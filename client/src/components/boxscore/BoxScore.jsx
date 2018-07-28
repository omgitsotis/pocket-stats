import React, {Component} from 'react';
import './boxscore.css'

class BoxScore extends Component {
    render() {
        const totals = this.props.totals;
        return(
            <div>
                <div className='row'>
                    <div className="col-lg-3">
                        <div className='card text-center stat-box'>
                            <h5 className='card-title'>Articles Add</h5>
                            <p className='card-text'>{totals.total_articles_added}</p>
                        </div>
                    </div>
                    <div className="col-lg-3">
                        <div className='card text-center stat-box'>
                            <h5 className='card-title'>Articles Read</h5>
                            <p className='card-text'>{totals.total_articles_read}</p>
                        </div>
                    </div>
                    <div className="col-lg-3">
                        <div className='card text-center stat-box'>
                            <h5 className='card-title'>Words add</h5>
                            <p className='card-text'>{totals.total_words_added}</p>
                        </div>
                    </div>
                    <div className="col-lg-3">
                        <div className='card text-center stat-box'>
                            <h5 className='card-title'>Words Read</h5>
                            <p className='card-text'>{totals.total_words_read}</p>
                        </div>
                    </div>
                    <div className="col-lg-3">
                        <div className='card text-center stat-box'>
                            <h5 className='card-title'>Time Read</h5>
                            <p className='card-text'>{totals.total_time_reading}</p>
                        </div>
                    </div>
                </div>
            </div>

        )
    }
}

export default BoxScore;
