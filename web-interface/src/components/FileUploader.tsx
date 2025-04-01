import React, { useEffect } from 'react';
import { Button } from '@heroui/button';
import { Spinner } from '@heroui/spinner';
import { Progress } from '@heroui/progress';

import { useFileContext } from '../contexts/FileContext.tsx';
import { useNotificationContext } from '../contexts/NotificationContext.tsx';

const FileUploader : React.FC = () => {
    const { files, filesReady, filesUploading, uploadFiles } = useFileContext();

    return (
        <>
        
        <div className="w-full px-2">
            <Button
                color="primary"
                isDisabled={!files || Object.keys(files).length === 0 || !filesReady}
                className="text-md w-full bg-primary-gradient"
                size="md"
                radius="sm"
                onPress={uploadFiles}
            >
                {((files && !filesReady) || filesUploading) ? <Spinner color="default"/> : "Send Files"}
            </Button>
        </div>
        </>
    )
}

export default FileUploader;