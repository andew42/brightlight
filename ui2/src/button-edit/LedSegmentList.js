import * as React from "react";
import './LedSegmentList.css';
import {Button, Icon, Image, Modal} from "semantic-ui-react";

export default class LedSegmentList extends React.Component {

    constructor(props) {
        super(props);

        this.state = {isOpen: true};
    }

    renderSegment(seg) {
        return <div className='ui image led-segment-list'>
            <Image label={seg.label} src={"/segment-icons/" + encodeURIComponent(seg.segment) + ".svg"}/>
            <div className='check-mark'>
                <Icon name={this.props.selectedItems.includes(seg.segment) ? 'checkmark box' : 'square outline'}/>
            </div>
        </div>;
    }

    render() {
        let segments = [
            {segment:'All', label:'All'}, {segment:'All Ceiling', label:'Ceiling'}, {segment:'All Wall', label:'Wall'},
            {segment:'Bedroom', label: 'Bedroom'}, {segment:'Bedroom Ceiling', label:'Ceiling'}, {segment:'Bedroom Wall', label:'Wall'},
            {segment:'Bathroom', label:'Bathroom'}, {segment:'Bathroom Ceiling', label:'Ceiling'}, {segment:'Bathroom Wall', label:'Wall'},
            {segment:'Chest', label:'Chest'}, {segment:'Chest Ceiling', label:'Ceiling'}, {segment:'Chest Wall', label:'Wall'},
            {segment:'Dressing', label:'Dressing'}, {segment:'Dressing Ceiling', label:'Ceiling'}, {segment:'Dressing Wall', label:'Wall'},
            {segment:'Curtains', label:'Curtains'}, {segment:'Door', label:'Door'}, {segment:'Light Switch', label:'Light Switch'}

        ];

        return <Modal trigger={this.props.trigger}
                      open={this.state.isOpen}
        >
            <Modal.Header>Select Light Segment</Modal.Header>
            <Modal.Content scrolling>
                <Modal.Description>
                    {segments.map(seg => this.renderSegment(seg))}
                </Modal.Description>
            </Modal.Content>
            <Modal.Actions>
                <Button primary content='OK'/>
                <Button content='Cancel' onClick={() => this.setState({isOpen: false})}/>
            </Modal.Actions>
        </Modal>
    }
}
