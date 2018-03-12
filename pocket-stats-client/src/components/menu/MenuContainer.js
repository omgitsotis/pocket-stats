import React, {Component} from 'react';

class MenuContainer extends Component {
    constructor(props) {
        super(props);
        this.state = {
            currentPage: 'menu'
        };
    }

    onButtonClick = (page) => {
        console.log(page);
        this.setState({currentPage: page});
    }

    render() {
        let section;
        switch (this.state.currentPage) {
            case 'init':
                section =
                    <InitSection
                        onBackClick={ () => this.onClick('menu')}
                        {...this.props}
                    />;
                break;
            case 'overview':
                section =
                    <div>
                        <button onClick={ () => this.props.onFetchDataClick() }>Fetch</button>
                    </div>
                break;
            case 'update':
                section =
                    <div>
                        <button onClick={ () => this.props.onUpdateClick() }>Update</button>
                    </div>
                break;
            default:
                section =
                    <div>
                        
                    </div>;
                break;
        }
        return (
            <div>
                {section}
            </div>
        )
    }
}
