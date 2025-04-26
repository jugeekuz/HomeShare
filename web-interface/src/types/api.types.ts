export interface SharingDetails {
    access:                 "rw" | "r" | "w";
    folder_name:            string;
    otp:                    string;
    expiration_date:        string;
}

export interface SharingGatewayDetails {
    link_url:   string;
    otp:       string;
}

export interface SharingResponse {
    link_url:           string;
    folder_id:          string;
}

export interface AuthSharingResponse {
    access_token:           string;
    folder_id:              string;
}


export interface TokenResponse {
    access_token:  string;
}
