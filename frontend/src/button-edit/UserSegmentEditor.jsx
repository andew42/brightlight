import React, {useState} from 'react';
import './UserSegmentEditor.css';
import {RangeEditor} from "./RangeEditor";

export function UserSegmentEditor(props) {
    const [baseSegment, setBaseSegment] = useState(props.predefinedSegments[0]);
    const [name, setName] = useState('');
    const [start, setStart] = useState(0);
    const [len, setLen] = useState(1);

    let editSeg = props.userSegments.find(x => x.name === name);

    return <div className='use-container'>
        <input className='use-input'
               placeholder='Name'
               value={name}
               onChange={e => {
                   setName(e.target.value);
                   let seg = props.userSegments.find(x => x.name === e.target.value);
                   if (seg) {
                       setBaseSegment(props.predefinedSegments.find(x => x.name === seg.base));
                       setStart(seg.start);
                       setLen(seg.length);
                   }
               }}/>
        <div className='use-based-on'>
            <span>Based on </span>
            <select value={baseSegment.name}
                    onChange={e => setBaseSegment(props.predefinedSegments.find(x => x.name === e.target.value))}>
                {props.predefinedSegments.map(s => <option key={s.name} value={s.name}>{s.name}</option>)}
            </select>
        </div>
        <RangeEditor label='Start' min={0} max={baseSegment.len} value={start} onPosChanged={p => setStart(p)}/>
        <RangeEditor label='Length' min={1} max={baseSegment.len} value={len} onPosChanged={p => setLen(p)}/>
        <button className='use-btn' onClick={() => console.info('Click')}>
            {editSeg ? `Edit ${name}` : `Add ${name}`}
        </button>
    </div>;
}
