import * as React from "react";
import './LedSegmentList.css';
import {Button, Checkbox, Image, Modal} from "semantic-ui-react";

export default class LedSegmentList extends React.Component {

    constructor(props) {
        super(props);

        this.state = {isOpen: true};
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
                    {this.props.segmentNames.map(seg => this.renderSegment(seg))}
                </Modal.Description>
            </Modal.Content>
            <Modal.Actions>
                <Button primary content='OK' onClick={this.props.onOk}/>
                <Button content='Cancel' onClick={this.props.onCancel}/>
            </Modal.Actions>
        </Modal>
    }
}
