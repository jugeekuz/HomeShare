import React, {useState} from "react";
import { Tabs, Tab, Card, CardBody } from "@heroui/react";
import { useNavigate } from "react-router-dom";
import UploadComponent from "../components/UploadComponent";
type Key = string | number;

const HomePage: React.FC = () => {
    const [selectedTab, setSelectedTab] = useState<Key>("upload")
    return (
        <div className="flex w-full h-full justify-center items-center">
            {/* Rectangle */}
            <Card className="max-w-full w-[400px]  bg-white">
                <CardBody className="justify-center items-center">
                    <Tabs
                        fullWidth
                        selectedKey={selectedTab}
                        size="sm"
                        aria-label="Tabs colors" 
                        radius="sm" 
                        variant="solid"
                        onSelectionChange={(key) => setSelectedTab(key)}
                        className="p-2"
                        classNames={{
                            base: "bg-white", 
                            tabList: "border border-gray-200 bg-white p-1 w-4/5 h-9 mx-auto -mb-1 ",
                            tab: "px-4 py-2 text-gray-500 font-normal cursor-pointer transition-all hover:bg-gray-100 data-[selected]:shadow-[0_4px_16px_0_rgba(0,0,0,0.12)] data-[selected]:font-medium data-[selected]:text-gray-800 text-[13px] h-[28px]"
                          }}
                    >
                        <Tab key="upload" title="Upload" className="w-full">
                            <UploadComponent/>
                        </Tab>
                        <Tab key="share" title="Share">

                        </Tab>
                    </Tabs>
                </CardBody>
            </Card>
        </div>
    );
};

export default HomePage;
