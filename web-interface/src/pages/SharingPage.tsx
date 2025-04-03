import React, {useState} from "react";
import { Tabs, Tab, Card, CardBody, Divider, DropdownMenu } from "@heroui/react";
import { useNavigate } from "react-router-dom";
import UploadComponent from "../components/UploadComponent";
import ShareComponent from "../components/ShareComponent";
import { FaLink } from "react-icons/fa6";
import { FileUploadProvider } from "../contexts/FileUploadContext.tsx";
import { LuUpload, LuDownload } from "react-icons/lu";
import DownloadComponent from "../components/DownloadComponent";
import { useSearchParams } from "react-router-dom";
type Key = string | number;

const SharingPage: React.FC = () => {
    const [selectedTab, setSelectedTab] = useState<Key>("download");
    const [searchParams] = useSearchParams();
    // 
    return (
        <div className="flex w-full h-full justify-center items-center">
            {/* Rectangle */}
            <FileUploadProvider>
            <Card className="max-w-full w-[440px]  bg-wprimary -mt-10">
                <CardBody className="justify-center items-center">
                    <Tabs
                        fullWidth
                        selectedKey={selectedTab}
                        size="sm"
                        aria-label="Tabs colors" 
                        radius="sm" 
                        variant="solid"
                        onSelectionChange={(key) => setSelectedTab(key)}
                        className=""
                        classNames={{
                            base: "bg-wprimary", 
                            tabList: "border border-gray-200 bg-wsecondary p-1 w-[70%] h-9 mx-auto mt-1",
                            tab: "px-4 py-2 text-gray-500 font-normal cursor-pointer transition-all hover:bg-gray-100 data-[selected]:shadow-[0_4px_16px_0_rgba(0,0,0,0.12)] data-[selected]:font-medium data-[selected]:text-gray-800 text-[13px] h-[28px]"
                          }}
                    >
                        <Tab key="download" title={
                            <div className="flex flex-row justify-center items-center ">
                                <LuDownload size={15} className='mr-2'/>
                                <span className="">Download</span>
                            </div>
                        } className="w-full">
                            <div className="-my-1 -mb-2">
                            <div className="flex flex-col justify-start items-center mt-1">
                                <div className="flex justify-center items-center">
                                    <span className="font-brsonoma font-normal text-gray-950 text-md mr-1">
                                        Sharing Folder <span className="bg-secondary-gradient bg-clip-text text-transparent font-bold">{searchParams.get("folder-name") || ""}</span>
                                    </span>
                                </div>
                                <span className="font-brsonoma font-light text-gray-500 text-xs">
                                    Explore what others shared or upload your files
                                </span>
                            </div>
                                <DownloadComponent/>
                            </div>
                        </Tab>
                        <Tab key="upload" title={
                            <div className="flex flex-row justify-center items-center ">
                                <LuUpload size={15} className='mr-2'/>
                                <span className="">Upload</span>
                            </div>
                        } className="w-full">
                            <div className="-my-1 -mb-2">
                                <div className="flex flex-col justify-start items-center mt-1">
                                    <div className="flex flex-col justify-center items-center mt-1 mb-1">
                                        <span className="font-brsonoma font-normal text-gray-950 text-md mr-1">
                                             Sharing Folder <span className="bg-secondary-gradient bg-clip-text text-transparent font-bold">{searchParams.get("folder-name") || ""}</span>
                                        </span>
                                        <span className="font-brsonoma font-light text-gray-500 text-xs">
                                            Explore what others shared or upload your files
                                        </span>
                                    </div>
                                </div>
                                <UploadComponent/>
                            </div>
                        </Tab>
                    </Tabs>
                </CardBody>
            </Card>
            </FileUploadProvider>
        </div>
    );
};

export default SharingPage;
