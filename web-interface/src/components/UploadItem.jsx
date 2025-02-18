import React, { useRef, useContext, useEffect, useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import SparkMD5 from 'spark-md5';
import { FileContext } from '../contexts/FileContext';

const UploadItem = ({file}) => {
    const { appendFilesMeta } = useContext(FileContext);

    const fileIdRef = useRef(uuidv4());
    const [fileMd5, setFileMd5] = useState(null);

    
    useEffect(() => {
        if (!file) return;

        const reader = new FileReader();
        let isCancelled = false;

        reader.onload = (e) => {
            if (isCancelled) return;
            const buffer = e.target.result;
            
            try {
                const hash = SparkMD5.ArrayBuffer.hash(buffer);
                setFileMd5(hash);
            } catch (error) {
                console.error('Error generating MD5 hash:', error);
            }
        };

        reader.onerror = (error) => {
            if (isCancelled) return;
            console.error('Error reading file:', error);
        };

        reader.readAsArrayBuffer(file);

        return () => {
            isCancelled = true;
            if (reader.readyState === FileReader.LOADING) {
                reader.abort();
                console.log('Aborted ongoing file read.');
            }
        };
    }, [file]);

    useEffect(() => {
        if (!fileMd5) return;

        const fileName = file.name;
        const parts = fileName.split(/(\..+)$/);

        const meta = {
            fileId:         fileIdRef.current,
            fileName:       parts[0],
            fileExtension:  parts[1] || '',
            md5Hash:        fileMd5
        }
        appendFilesMeta(meta);
    }, [file, fileMd5])

    return (
        <div>File</div>
    )
}

export default UploadItem