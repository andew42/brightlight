import * as React from "react";
import {Fragment} from "react";
import NameEditor from "./NameEditor";
import './ButtonEdit.css';
import {saveButtons} from "../server-proxy/buttons";
import {Button} from "semantic-ui-react";
import LedSegment from "./LedSegment";
import LedSegmentList from "./LedSegmentList";

// Shows a list of editors that change to suit the animation
export default class ButtonEdit extends React.Component {
    constructor(props) {
        super(props);

        // TODO
        this.segmentList = ['Static', 'Twinkle', 'Rainbow'];

        // We expect to be passed all buttons and button to edit in location state
        this.buttons = props.history.location.state.buttons;
        this.button = props.history.location.state.button;

        this.state = {
            // Initial state of button we are editing
            name: this.button.name,
            segments: JSON.parse(JSON.stringify(this.button.segments)),
// TODO            selectedSegments: this.state.segments.map(s => s.segment)
        };
        this.onOK = this.onOK.bind(this);
        this.onAddSegment = this.onAddSegment.bind(this);
    }

    onOK() {
        // Update button object and save to server
        this.button.name = this.state.name;
        // TODO HANDLE ERRORS BETTER
        saveButtons(this.buttons, () => this.props.history.goBack(), e => console.error(e));
    }

    onRemoveSegment(segment) {
        console.info('removing segment ' + segment.segment);
        this.setState({segments: this.state.segments.filter(seg => seg !== segment)})
    }

    onAddSegment() {
        console.info('adding segment');
        // TODO
    }

    toggleSelectedSegment(segName) {
        // TODO
        // if (state.selectedSegments.includes(segName))
        //     this.setState({selectedSegments: state.selectedSegments.filter(n => n != segName)});
        // else
        //     this.setState({selectedSegments: state.s})
    }

    render() {
        let selectedSegments = this.state.segments.map(s => s.segment);
        let key = 1;
        return <div className="button-edit-editor-list">
            <Fragment>
                <NameEditor name={this.state.name} onNameChanged={name => this.setState({name: name})}/>

                {this.state.segments.map(segment => (
                    <LedSegment key={key++} segment={segment} onRemove={seg => this.onRemoveSegment(seg)}/>))}

                <div className='ok-cancel-container'>
                    <LedSegmentList selectedItems={selectedSegments}
                                    toggleSelectedSegment={seg => console.info(seg)}
                                    trigger={<Button icon='plus' circular floated='left'></Button>}/>
                    <Button primary onClick={this.onOK}>OK</Button>
                    <Button secondary onClick={() => this.props.history.goBack()}>Cancel</Button>
                </div>
            </Fragment>
        </div>
    }
}
