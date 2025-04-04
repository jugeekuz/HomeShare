import React, { useEffect, useState } from 'react';
import { FileDownloadItem, FileIconType } from '../types';
import { FileIcon, defaultStyles, DefaultExtensionType } from 'react-file-icon';

import { useFileDownloadContext } from '../contexts/FileDownloadContext';
import { useNotificationContext } from '../contexts/NotificationContext';

const DownloadItem : React.FC<{fileName: string}> = ({ fileName }) => {
    const [fileStyle, setFileStyle] = useState<FileIconType>({fileExtension: '', fileStyle: defaultStyles});
    const { notifyError } = useNotificationContext();
    const {files} = useFileDownloadContext();
    const [file, setFile] = useState<FileDownloadItem | null>(null);

    useEffect(() => {
        if (!files || !fileName) return;
        
        if (!(fileName in files)) {
            notifyError("File Error", `File name ${fileName} wasn't found in the store`);
            return;
        }
        setFile(files[fileName]);        
    }, []);

    useEffect(() => {
        if (!file || !file?.fileExtension) return;
        const fileExt = file.fileExtension.replace(/^\./, "").toLowerCase() as DefaultExtensionType;

        setFileStyle({
            fileExtension: fileExt,
            fileStyle: (defaultStyles[fileExt] || defaultStyles)
        })
    }, [file])
    
    return (
        <FileIcon color='#728eab' extension={fileStyle.fileExtension} {...fileStyle.fileStyle}  />
    )
}

export default DownloadItem;