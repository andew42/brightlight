import * as React from "react";
import './LedSegmentChooser.css';
import Dialog from "../dialog/Dialog";
import {UserSegmentEditor} from "./UserSegmentEditor";

export default class LedSegmentChooser extends React.Component {

    constructor(props) {
        super(props);
        this.state = {open: false};
    }

    render() {
        return <>
            <span onClick={() => this.setState({open: true})}>{this.props.trigger}</span>
            <Dialog open={this.state.open}
                    header='Select Light Segment'
                    actions={[
                        {key: 'ok', content: 'OK', primary: true, onClick: () => {
                            this.props.onOk && this.props.onOk();
                            this.setState({open: false});
                        }},
                        {key: 'cancel', content: 'Cancel', onClick: () => {
                            this.props.onCancel();
                            this.setState({open: false});
                        }}
                    ]}
                    onClose={() => this.setState({open: false})}>
                <div>
                    {this.props.allSegments.map(seg => this.renderSegment(seg))}
                </div>
                <div>
                    {this.props.userSegments.map(seg => this.renderUserSegment(seg))}
                </div>
                <hr className='lsc-divider'/>
                <div className='lsc-divider-label'>New User Segment</div>
                <UserSegmentEditor predefinedSegments={this.props.allSegments}
                                   userSegments={this.props.userSegments}/>
            </Dialog>
        </>;
    }

    renderSegment(seg) {
        const checked = this.props.checkedSegmentNames.includes(seg.name);
        return (
            <label className='lsc-led-segment-list-item' key={seg.name}>
                <img className='lsc-icon'
                     alt=''
                     src={"/segment-icons/" + encodeURIComponent(seg.name) + ".svg"}
                     onError={e => e.target.style.display = 'none'}/>
                <span className='lsc-segment-label'>{seg.label}</span>
                <input type='checkbox'
                       className='lsc-check'
                       checked={checked}
                       onChange={() => this.props.toggleCheckedSegment(seg)}/>
            </label>
        );
    }

    renderUserSegment(seg) {
        return (
            <div className='lsc-led-user-segment-list-item' key={seg.name}>
                <label className='lsc-user-label'>
                    <input type='checkbox'
                           checked={this.props.checkedSegmentNames.includes(seg.name)}
                           onChange={() => this.props.toggleCheckedSegment(seg)}/>
                    {seg.name}
                </label>
            </div>
        );
    }
}
