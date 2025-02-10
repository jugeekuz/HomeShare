import React, { useState, useEffect, useContext, use } from 'react';
import config from '../configs/config';
import { Button } from '@heroui/button';
import { Progress } from '@heroui/progress';
import { uploadChunk, createChunk, chunkAndUpload } from '../services/chunkFile';

import { FileContext } from '../contexts/FileContext';

const FileUploader = () => {
  const { files, setFiles, progress, setProgress } = useContext(FileContext);
  const [fileProgresses, setFileProgresses] = useState([]);
  const uploadFile = async () => {
    if (!files) return;
    chunkAndUpload(setProgress, files);
  } 
  useEffect(() => {
    if (fileProgresses.length === 0) {
      setProgress(0);
      return;
    }
    const total = fileProgresses.reduce((acc, curr) => acc + curr, 0) / fileProgresses.length;
    setProgress(total);
  }, [fileProgresses, setProgress]);

  const uploadFiles = async () => {
    if (!files || files.length === 0) return;

    setFileProgresses(Array(files.length).fill(0));

    for (let i = 0; i < files.length; i++) {
      const currentFile = files[i];
      await chunkAndUpload(
        (progress) => {
          setFileProgresses(prev => {
            const newProgresses = [...prev];
            newProgresses[i] = progress;
            return newProgresses;
          });
        },
        currentFile
      );
    }

  };

  useEffect(() => {
    if (progress !== 100) return;

    setFiles([]);
    setProgress(0)
  }, [progress])

  return (
    <>
      { progress !== 0 && (
        <Progress 
          isStriped 
          aria-label="Loading..." 
          className="w-[70%]" 
          color="secondary" 
          value={progress}
        />
      )}
      <Button
        color="primary"
        isDisabled={!files || files.length === 0}
        className="text-md w-[80%]"
        size="lg"
        onPress={uploadFiles}
      >
        Send Files
      </Button>
    </>
  )
}

export default FileUploader;