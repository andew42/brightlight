import * as React from "react";
import {Fragment} from "react";
import './ButtonEditor.css';
import {saveButtons} from "../server-proxy/buttons";
import {Button} from "semantic-ui-react";
import LedSegmentChooser from "./LedSegmentChooser";
import {getStaticData} from "../server-proxy/staticData";
import {LedSegmentEditor} from "./LedSegmentEditor";
import {NameEditor} from "./NameEditor";

// Shows a list of editors that change to suit the animation
export default class ButtonEditor extends React.Component {
    constructor(props) {
        super(props);

        // We expect to be passed all buttons and button to edit in location state
        this.buttons = props.history.location.state.buttons;
        this.button = props.history.location.state.button;

        // Copy the initial state of the button we are editing e.g.
        // {
        //     "name": "Cylon",
        //     "segments": [
        //     {
        //         "segment": "Curtains",
        //         "animation": "Cylon"
        //     }]
        // }
        let buttonSegmentsCopy = JSON.parse(JSON.stringify(this.button.segments));
        this.state = {
            // Static data
            allSegmentNames: [],
            allAnimationNames: [],

            // The button we are editing
            name: this.button.name,
            segments: buttonSegmentsCopy,
            selectedSegments: buttonSegmentsCopy.map(s => s.segment),
        };

        // Callback bindings for this
        this.onOK = this.onOK.bind(this);
        this.toggleSelectedSegment = this.toggleSelectedSegment.bind(this);
        this.editSegmentListCancel = this.editSegmentListCancel.bind(this);
        this.editSegmentListOk = this.editSegmentListOk.bind(this);
    }

    // Get static data from server when mounting
    componentDidMount() {
        getStaticData(sd => this.setState({allSegmentNames: sd.segments, allAnimationNames: sd.animations}),
            (xhr) => console.error(xhr)); // TODO: Report errors to user
    }

    onOK() {
        // Update button object and save to server
        this.button.name = this.state.name;
        // TODO: Full update and better error handling
        saveButtons(this.buttons, () => this.props.history.goBack(),
            e => console.error(e));
    }

    onRemoveSegment(segment) {
        console.info('removing segment ' + segment.segment);
        this.setState({
            segments: this.state.segments.filter(seg => seg !== segment),
            selectedSegments: this.state.selectedSegments.filter(seg => seg !== segment.segment)
        })
    }

    toggleSelectedSegment(segName) {
        console.info('Toggle:' + segName);
        if (this.state.selectedSegments.includes(segName))
            this.setState({selectedSegments: this.state.selectedSegments.filter(n => n !== segName)});
        else {
            let selectedSegmentCopy = JSON.parse(JSON.stringify(this.state.selectedSegments));
            selectedSegmentCopy.push(segName);
            this.setState({selectedSegments: selectedSegmentCopy});
        }
    }

    editSegmentListCancel() {
        this.setState({
            // Restore selected segment list as user canceled
            selectedSegments: this.state.segments.map(s => s.segment)
        });
    }

    editSegmentListOk() {
        let segments = this.state.segments;
        let selectedSegments = this.state.selectedSegments;
        // Build new segments from selectedSegments
        let updatedSegments = selectedSegments.map(seg => {
            // Use existing segment if it's still selected
            let foundSegment = segments.find(s => s.segment === seg);
            if (foundSegment !== undefined)
                return foundSegment;
            // Here if we have a new segment to add (default static animation)
            return {
                "segment": seg,
                "animation": "Static",
                "params": "3f3f3f"
            }
        });
        this.setState({segments: updatedSegments});
    }

    // The animation on a particular segment has changed
    onAnimationNameChange(seg, val) {

        let newSegments = this.state.segments.map(s => s !== seg ? s : Object.assign({}, s, {animation: val}));
        this.setState({segments: newSegments});
    }

    render() {
        let key = 1;
        let animationNames = this.state.allAnimationNames.map(n => ({'text': n, 'value': n}));
        return <div className="button-edit-editor-list">
            <Fragment>
                <NameEditor name={this.state.name} onNameChanged={name => this.setState({name: name})}/>

                {this.state.segments.map(segment => (
                    <LedSegmentEditor key={key++}
                                      segment={segment}
                                      onRemove={seg => this.onRemoveSegment(seg)}
                                      onAnimationNameChange={(e, d) => this.onAnimationNameChange(segment, d.value)}
                                      animationNames={animationNames}/>))}

                <div className='ok-cancel-container'>
                    <LedSegmentChooser allSegmentNames={this.state.allSegmentNames}
                                       onCancel={this.editSegmentListCancel}
                                       onOk={this.editSegmentListOk}
                                       checkedSegmentNames={this.state.selectedSegments}
                                       toggleCheckedSegment={seg => this.toggleSelectedSegment(seg)}
                                       trigger={<Button icon='plus'
                                                        circular
                                                        floated='left'/>}/>
                    <Button primary onClick={this.onOK}>OK</Button>
                    <Button secondary onClick={() => this.props.history.goBack()}>Cancel</Button>
                </div>
            </Fragment>
        </div>
    }
}
