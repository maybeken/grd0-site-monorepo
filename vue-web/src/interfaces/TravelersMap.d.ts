export interface MapLocation {
    title: string,
    icon: string,
    color: string,
    pos: [number, number], // Record coordinate by North-East
    subtitle?: string,
    displayAt?: number,
    hideAt?: number,
    textColor?: string,
};