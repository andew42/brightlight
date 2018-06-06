import React, {Component} from 'react';
import {Helmet} from "react-helmet";
import './App.css';
import ButtonPad from "./buttons/ButtonPad";
import {BrowserRouter, Link, Route, Switch} from 'react-router-dom'
import ButtonEditor from "./button-edit/ButtonEditor";
import {getStaticData} from "./server-proxy/staticData";
import {getButtons} from "./server-proxy/buttons";
import {runAnimation} from "./server-proxy/animation";
import Virtual from "./virtual/Virtual";

// Home page is just a bunch of links for now
const Home = () => (
    <nav>
        <Link to="/buttons">Buttons</Link>
        <hr/>
        <Link to="/virtual">Virtual</Link>
    </nav>
);

// The main application navigates between components
export default class App extends Component {

    constructor(props) {
        super(props);
        this.state = {allButtons: [], allSegments: [], allAnimations: []}
    }

    // Get static data from server when mounting
    componentDidMount() {
        getButtons(
            buttons => this.setState({allButtons: buttons}),
            (xhr) => console.error(xhr)); // TODO: Report errors to user

        getStaticData(
            sd => this.setState({allSegments: sd.segments, allAnimations: sd.animations}),
            (xhr) => console.error(xhr)); // TODO: Report errors to user
    }

    onButtonChanged(button) {
        let allButtons = this.state.allButtons.map(b => b.key === button.key ? button : b);
        this.setState((props, state) => {
            return {...state, allButtons: allButtons}
        });
    }

    onSaveButtonEdit() {
        console.info('onSaveButtonEdit');
        // TODO SAVE!
    }

    findButton(key) {
        return this.state.allButtons.find(x => x.key === key);
    }

    onButtonTap(key) {
        console.info('onButtonTap ' + key);
        runAnimation(this.findButton(key).segments);
    }

    onButtonPress(history, key) {
        history.push('/button-edit', {'buttonKey': key});
    }

    render() {

        let ButtonPadWithProps = props => {
            return (<ButtonPad allButtons={this.state.allButtons}
                               onButtonTap={key => this.onButtonTap(key)}
                               onButtonPress={(history, key) => this.onButtonPress(history, key)}
                               {...props}/>);
        };

        let allButtons = this.state.allButtons;
        let ButtonEditorWithProps = props => {
            return (<ButtonEditor allButtons={allButtons}
                                  allAnimations={this.state.allAnimations}
                                  allSegments={this.state.allSegments}
                                  onButtonChanged={button => this.onButtonChanged(button)}
                                  onOk={() => this.onSaveButtonEdit()}
                                  {...props}/>);
        };

        return (
            <div className="App" onContextMenu={e => e.preventDefault()}>
                <Helmet>
                    <meta name="apple-mobile-web-app-capable" content="yes"/>
                    <title>Bright Light</title>
                </Helmet>
                <BrowserRouter>
                    <Switch>
                        <Route path="/buttons" render={ButtonPadWithProps}/>
                        <Route path="/button-edit" render={ButtonEditorWithProps}/>
                        <Route path="/virtual" component={Virtual}/>
                        <Route path="/" component={Home}/>
                    </Switch>
                </BrowserRouter>
            </div>
        );
    }
}
