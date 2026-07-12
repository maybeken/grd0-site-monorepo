export interface Columns {
  [key: string]: {
    display_name: string | Function
    formatter: Function
  }
}

export interface Action {
  name: string
  display_name: string
  data_key: string
}
