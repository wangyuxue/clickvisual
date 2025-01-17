import queryResultStyles from "@/pages/DataLogs/components/QueryResult/index.less";
import RawLogsIndexes from "@/pages/DataLogs/components/QueryResult/Content/RawLog/RawLogsIndexes";
import { Spin } from "antd";
import classNames from "classnames";
import HighCharts from "@/pages/DataLogs/components/QueryResult/Content/RawLog/HighCharts";
import RawLogs from "@/pages/DataLogs/components/QueryResult/Content/RawLog/RawLogs";
import { useModel } from "@@/plugin-model/useModel";
import { useIntl } from "umi";

const RawLogContent = () => {
  const { logsLoading, highChartLoading, isHiddenHighChart } =
    useModel("dataLogs");

  const i18n = useIntl();

  return (
    <div className={queryResultStyles.content}>
      <RawLogsIndexes />
      <div className={queryResultStyles.queryDetail}>
        <Spin
          spinning={highChartLoading}
          tip={i18n.formatMessage({ id: "spin" })}
          wrapperClassName={classNames(
            queryResultStyles.querySpinning,
            isHiddenHighChart
              ? queryResultStyles.highChartsHidden
              : queryResultStyles.highCharts
          )}
        >
          <HighCharts />
        </Spin>
        <Spin
          spinning={logsLoading}
          tip={i18n.formatMessage({ id: "spin" })}
          wrapperClassName={classNames(
            queryResultStyles.querySpinning,
            queryResultStyles.logs
          )}
        >
          <RawLogs />
        </Spin>
      </div>
    </div>
  );
};
export default RawLogContent;
