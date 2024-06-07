import { RawTimeRange } from '@grafana/data';
import { DataQuery } from '@grafana/schema';
export const TOP_BAR_LEVEL_HEIGHT = 40;

export interface HistoryEntryAppView {
  name: string;
  url: string;
  timeRange?: RawTimeRange;
  query?: DataQuery;
}

export interface HistoryEntryApp {
  name: string;
  views: HistoryEntryAppView[];
}