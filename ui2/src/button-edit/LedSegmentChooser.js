import * as React from "react";
import './LedSegmentChooser.css';
import {Button, Checkbox, Image, Modal} from "semantic-ui-react";

export default class LedSegmentChooser extends React.Component {

    // props
    //   selectedItems
    //   trigger - button used to open chooser
    //   isOpen - is the editor visible
    //   allSegmentNames - the list of segments to choose from
    //   toggleSelectedSegment - function to call when a segment is toggled
    //   onOK - function to call when user hits OK
    //   onCancel - function to call when user hits cancel
    constructor(props) {
        super(props);
    }

    renderSegment(seg) {
        return <div className='ui image led-segment-list' key={seg.segment}>
            <Image label={seg.label} src={"/segment-icons/" + encodeURIComponent(seg.segment) + ".svg"}/>
            <div className='check-mark'>
                <Checkbox checked={this.props.selectedItems.includes(seg.segment)}
                          onChange={() => this.props.toggleSelectedSegment(seg.segment)}/>
            </div>
        </div>;
    }

    render() {
        return <Modal trigger={this.props.trigger} open={this.props.isOpen}>
            <Modal.Header>Select Light Segment</Modal.Header>
            <Modal.Content scrolling>
                <Modal.Description>
                    {this.props.allSegmentNames.map(seg => this.renderSegment(seg))}
                </Modal.Description>
            </Modal.Content>
            <Modal.Actions>
                <Button primary content='OK' onClick={this.props.onOk}/>
                <Button content='Cancel' onClick={this.props.onCancel}/>
            </Modal.Actions>
        </Modal>
    }
}
