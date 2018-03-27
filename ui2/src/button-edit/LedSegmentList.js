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
        let segments = [
            {segment: 'All', label: 'All'},
            {segment: 'All Ceiling', label: 'Ceiling'},
            {segment: 'All Wall', label: 'Wall'},
            {segment: 'Bedroom', label: 'Bedroom'},
            {segment: 'Bedroom Ceiling', label: 'Ceiling'},
            {segment: 'Bedroom Wall', label: 'Wall'},
            {segment: 'Bathroom', label: 'Bathroom'},
            {segment: 'Bathroom Ceiling', label: 'Ceiling'},
            {segment: 'Bathroom Wall', label: 'Wall'},
            {segment: 'Chest', label: 'Chest'},
            {segment: 'Chest Ceiling', label: 'Ceiling'},
            {segment: 'Chest Wall', label: 'Wall'},
            {segment: 'Dressing', label: 'Dressing'},
            {segment: 'Dressing Ceiling', label: 'Ceiling'},
            {segment: 'Dressing Wall', label: 'Wall'},
            {segment: 'Curtains', label: 'Curtains'},
            {segment: 'Door', label: 'Door'},
            {segment: 'Light Switch', label: 'Light Switch'}
        ];

        return <Modal trigger={this.props.trigger} open={this.props.isOpen}>
            <Modal.Header>Select Light Segment</Modal.Header>
            <Modal.Content scrolling>
                <Modal.Description>
                    {segments.map(seg => this.renderSegment(seg))}
                </Modal.Description>
            </Modal.Content>
            <Modal.Actions>
                <Button primary content='OK' onClick={this.props.onOk}/>
                <Button content='Cancel' onClick={this.props.onCancel}/>
            </Modal.Actions>
        </Modal>
    }
}
