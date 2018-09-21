import React, {Component} from 'react';
import {Helmet} from "react-helmet";
import './App.css';
import ButtonPad from "./buttons/ButtonPad";
import {BrowserRouter, Link, Route, Switch} from 'react-router-dom'
import ButtonEditor from "./button-edit/ButtonEditor";
import {getStaticData} from "./server-proxy/staticData";
import {getButtons, saveButtons} from "./server-proxy/buttons";
import {runAnimation} from "./server-proxy/animation";
import Virtual from "./virtual/Virtual";
import {OpenWebSocket} from "./server-proxy/webSocket";

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
        this.state =
            {
                allButtons: undefined,
                allSegments: undefined,
                allAnimations: undefined,
                activeButtonKey: 0,
                buttonPadVersion: 1
            };
    }

    // Get static data from server when mounting
    componentDidMount() {

        console.info('getting static data');
        getStaticData(
            sd => this.setState({allSegments: sd.segments, allAnimations: sd.animations}),
            (xhr) => console.error(xhr));

        // Subscribe to button state updates (we get immediate update on connection)
        this.ws = OpenWebSocket('ButtonState', bs => {
            console.info('button state changed to ', bs);

            if (this.state.activeButtonKey !== bs.ActiveButtonKey)
                this.setState({activeButtonKey: bs.ActiveButtonKey});

            if (this.state.buttonPadVersion !== bs.ButtonPadVersion) {
                this.setState({buttonPadVersion: bs.ButtonPadVersion});
                this.getButtonConfig();
            }
        });
    }

    getButtonConfig() {
        console.info('getting button configuration');
        getButtons(
            buttons => this.setState({allButtons: buttons}),
            (xhr) => console.error(xhr));
    }

    // Close web socket
    componentWillUnmount() {
        if (this.ws !== undefined)
            this.ws.close()
    }

    onButtonChanged(button) {
        let allButtons = this.state.allButtons.map(b => b.key === button.key ? button : b);
        this.setState((props, state) => {
            return {...state, allButtons: allButtons}
        });
        runAnimation(button);
    }

    onSaveButtonEdit() {
        console.info('onSaveButtonEdit');
        saveButtons(
            this.state.allButtons,
            () => console.info('button state saved'),
            (xhr) => console.error(xhr)); // TODO: Report errors to user
    }

    findButton(key) {
        return this.state.allButtons.find(x => x.key === key);
    }

    onButtonTap(key) {
        runAnimation(this.findButton(key));
    }

    onButtonPress(history, key) {
        runAnimation(this.findButton(key));
        history.push('/button-edit', {'buttonKey': key});
    }

    render() {

        let ButtonPadWithProps = props => {
            return (<ButtonPad allButtons={this.state.allButtons}
                               activeButtonKey={this.state.activeButtonKey}
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
