import React from "react";
import { FileProvider } from "../contexts/FileContext.tsx";
import FileLoader from "../components/FileLoader.tsx";
import FileList from "../components/FileList.tsx";
import FileUploader from "../components/FileUploader.tsx";

const UploadPage: React.FC = () => {
  return (
    <div className="flex h-full w-full justify-center items-center px-1">
        <div className="flex flex-col justify-center items-center w-[25rem] xl:w-[32rem] bg-white rounded-xl shadow-lg -mt-20">
            <FileProvider>
                <div className="flex flex-col items-center justify-center w-full h-64">
                    <FileLoader />
                </div>
                <FileList />
                <div className="flex flex-col justify-center items-center w-full">
                    <FileUploader />
                </div>
            </FileProvider>
        </div>
    </div>
  );
};

export default UploadPage;
