import * as React from "react";
import './LedSegmentChooser.css';
import {Checkbox, Divider, Image, Modal} from "semantic-ui-react";
import {UserSegmentEditor} from "./UserSegmentEditor";

export default class LedSegmentChooser extends React.Component {

    render() {
        return <Modal trigger={this.props.trigger}
                      header='Select Light Segment'
                      content={<div className='scrolling content'>
                          <div className='description'>
                              {this.props.allSegments.map(seg => this.renderSegment(seg))}
                          </div>
                          <div>
                              {this.props.userSegments.map(seg => this.renderUserSegment(seg))}
                          </div>
                          <Divider horizontal>New User Segment</Divider>
                          <UserSegmentEditor predefinedSegments={this.props.allSegments}
                                             userSegments={this.props.userSegments}/>
                      </div>}
                      actions={[
                          {key: 'ok', content: 'OK', primary: true},
                          {key: 'cancel', content: 'Cancel'}
                      ]}
                      onActionClick={e => e.target.textContent === 'OK' ?
                          (this.props.onOk !== undefined && this.props.onOk()) :
                          this.props.onCancel()}>
        </Modal>
    }

    renderSegment(seg) {
        if (!seg.icon)
            return <div className='ui image lsc-led-segment-list-item-no-icon'
                        key={seg.name}>
                <Image label={seg.label}/>
                <div className='lsc-check-mark-no-icon'>
                    <Checkbox checked={this.props.checkedSegmentNames.includes(seg.name)}
                              onChange={() => this.props.toggleCheckedSegment(seg)}/>
                </div>
            </div>;

        return <div className='ui image lsc-led-segment-list-item'
                    key={seg.name}>
            <Image label={seg.label}
                   src={"/segment-icons/" + encodeURIComponent(seg.name) + ".svg"}/>
            <div className='lsc-check-mark'>
                <Checkbox checked={this.props.checkedSegmentNames.includes(seg.name)}
                          onChange={() => this.props.toggleCheckedSegment(seg)}/>
            </div>
        </div>;
    }

    renderUserSegment(seg) {
        return <div className='lsc-led-user-segment-list-item'
                    key={seg.name}>
            <Checkbox label={seg.name}
                      checked={this.props.checkedSegmentNames.includes(seg.name)}
                      onChange={() => this.props.toggleCheckedSegment(seg)}/>
        </div>;
    }
}
