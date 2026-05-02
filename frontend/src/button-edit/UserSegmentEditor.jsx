import React, {useState} from 'react';
import './UserSegmentEditor.css';
import {Button, Dropdown, Header, Input} from "semantic-ui-react";
import {RangeEditor} from "./RangeEditor";

export function UserSegmentEditor(props) {
    const [baseSegment, setBaseSegment] = useState(props.predefinedSegments[0]);
    const [name, setName] = useState('');
    const [start, setStart] = useState(0);
    const [len, setLen] = useState(1);

    let options = props.predefinedSegments.map(s =>
        ({key: s.name, text: s.name, value: s.name, content: s.name}));

    let editSeg = props.userSegments.find(x => x.name === name);

    return <div className='use-container'>
        <Input fluid placeholder='Name' onChange={(e, d) => {
            setName(d.value);
            let editSeg = props.userSegments.find(x => x.name === d.value);
            if (editSeg) {
                setBaseSegment(props.predefinedSegments.find(x => x.name === editSeg.base));
                setStart(editSeg.start);
                setLen(editSeg.length);
            }
        }}/>
        <Header as='h4'>
            <Header.Content>
                Based on{' '}
                <Dropdown
                    inline
                    options={options}
                    value={baseSegment.name}
                    onChange={(_, d) => {
                        setBaseSegment(props.predefinedSegments.find(x => x.name === d.value));
                    }}
                />
            </Header.Content>
        </Header>
        <div>
            <RangeEditor label='Start' min={0} max={baseSegment.len} value={start} onPosChanged={p => setStart(p)}/>
            <RangeEditor label='Length' min={1} max={baseSegment.len} value={len} onPosChanged={p => setLen(p)}/>
        </div>
        {editSeg
            ? (<Button onClick={() => console.info('Click')}>Edit {name}</Button>)
            : (<Button onClick={() => console.info('Click')}>Add {name}</Button>)
        }
    </div>
}
