import alarmStyles from "@/pages/Alarm/Rules/styles/index.less";
import { Divider, message, Space, Table, Tooltip } from "antd";
import { ColumnsType } from "antd/es/table";
import { useIntl } from "umi";
import { useModel } from "@@/plugin-model/useModel";
import { useEffect } from "react";
import { AlarmType } from "@/services/alarm";
import IconFont from "@/components/IconFont";
import { EditOutlined, FileTextOutlined } from "@ant-design/icons";
import DeletedModal from "@/components/DeletedModal";
import classNames from "classnames";
import { useDebounceFn } from "ahooks";
import useAlarmEnums from "@/pages/Alarm/hooks/useAlarmEnums";
import { ALARM_HISTORY_PATH } from "@/models/alarms/useAlarmHistory";
import useUrlState from "@ahooksjs/use-url-state";
import useTimeUnits from "@/hooks/useTimeUnits";

import moment from "moment";

const AlarmTable = () => {
  const [urlState] = useUrlState<any>();
  const i18n = useIntl();
  const { FixedTimeUnits } = useTimeUnits();

  const {
    alarmList,
    operations,
    alarmDraw,
    doGetAlarms,
    doDeletedAlarm,
    currentPagination,
    onChangeRowAlarm,
    onChangePagination,
  } = useModel("alarm");
  const { AlarmStatus } = useAlarmEnums();

  const searchQuery = {
    name: operations.inputName,
    did: operations.selectDid,
    tid: operations.selectTid,
    status: operations.statusId,
    ...currentPagination,
  };

  const handleEdit = (record: AlarmType) => {
    alarmDraw.onChangeIsEditor(true);
    onChangeRowAlarm(record);
    alarmDraw.onChangeVisibleDraw(true);
  };

  const handleInfo = (id: number) => {
    alarmDraw.doGetAlarmInfo.run(id).then((res) => {
      if (res?.code !== 0) return;
      alarmDraw.setAlarmInfo(res.data);
      alarmDraw.onChangeVisibleInfo(true);
    });
  };

  const handleDelete = (record: AlarmType) => {
    DeletedModal({
      content: i18n.formatMessage(
        { id: "alarm.rules.deleted.content" },
        { alarm: record.alarmName }
      ),
      onOk: () => {
        const hideMessage = message.loading(
          {
            content: i18n.formatMessage({ id: "alarm.rules.deleted.loading" }),
            key: "alarm",
          },
          0
        );
        doDeletedAlarm
          .run(record.id)
          .then((res) => {
            if (res?.code !== 0) {
              hideMessage();
              return;
            }
            doGetAlarms.run(searchQuery);
            message.success(
              {
                content: i18n.formatMessage({ id: "alarm.rules.deleted" }),
                key: "alarm",
              },
              3
            );
          })
          .catch(() => hideMessage());
      },
    });
  };

  const doUpdateStatus = useDebounceFn(
    (alarm: AlarmType) => {
      const isOpen = alarm.status !== 1;
      const hideMessage = message.loading(
        {
          content: i18n.formatMessage(
            {
              id: `alarm.rules.${isOpen ? "close" : "open"}.loading`,
            },
            { alarmName: alarm.alarmName }
          ),
          key: "status",
        },
        0
      );
      alarmDraw.doUpdatedAlarm
        .run(alarm.id, { status: isOpen ? 1 : 2 })
        .then((res) => {
          if (res?.code !== 0) {
            hideMessage();
            return;
          }
          message.success(
            {
              content: i18n.formatMessage(
                { id: `alarm.rules.${isOpen ? "close" : "open"}.success` },
                { alarmName: alarm.alarmName }
              ),
              key: "status",
            },
            3
          );
          doGetAlarms.run(searchQuery);
        })
        .catch(() => hideMessage());
    },
    { wait: 500 }
  ).run;

  useEffect(() => {
    if (urlState && urlState.name) {
      searchQuery.name = urlState.name;
      if (urlState.did || urlState.status) {
        const data = {
          ...searchQuery,
          ...urlState,
        };
        doGetAlarms.run(data);
        searchQuery.did = urlState.did * 1 || undefined;
        searchQuery.tid = urlState.tid * 1 || undefined;
        searchQuery.status = urlState.status * 1 || undefined;
      } else {
        doGetAlarms.run(searchQuery);
      }
    } else {
      doGetAlarms.run(searchQuery);
    }
  }, []);

  const column: ColumnsType<any> = [
    {
      title: i18n.formatMessage({ id: "alarm.rules.table.alarmName" }),
      dataIndex: "alarmName",
      align: "center",
      ellipsis: { showTitle: true },
      render: (alarmName: string, record: AlarmType) => (
        <Tooltip title={alarmName}>
          <div
            style={{ color: "red" }}
            className={classNames(
              alarmStyles.columnsEllipsis,
              alarmStyles.columnsTag
            )}
          >
            <a
              href={`${ALARM_HISTORY_PATH}${record.id}`}
              target={`alarm-${record.id}`}
            >
              <span>{alarmName}</span>
            </a>
          </div>
        </Tooltip>
      ),
    },
    {
      title: i18n.formatMessage({ id: "alarm.rules.inspectionFrequency" }),
      dataIndex: "interval",
      align: "center",
      render: (interval: number, record: AlarmType) => {
        const unit = FixedTimeUnits.filter(
          (item) => item.key === record.unit
        )[0];
        return (
          <span>
            {interval}&nbsp;{unit.label}
          </span>
        );
      },
    },
    {
      title: i18n.formatMessage({ id: "user" }),
      dataIndex: "user",
      ellipsis: { showTitle: true },
      render: (user: any) => (
        <Tooltip title={user.nickname}>
          <div className={alarmStyles.columnsEllipsis}>
            <span>{user.username}</span>
          </div>
        </Tooltip>
      ),
    },
    {
      title: i18n.formatMessage({ id: "description" }),
      dataIndex: "desc",
      width: "40%",
      ellipsis: { showTitle: true },
      render: (desc: string) => (
        <Tooltip title={desc}>
          <div className={alarmStyles.columnsEllipsis}>
            <span>{desc}</span>
          </div>
        </Tooltip>
      ),
    },
    {
      title: i18n.formatMessage({ id: "utime" }),
      dataIndex: "utime",
      ellipsis: { showTitle: true },
      render: (utime: any) => (
        <div className={alarmStyles.columnsEllipsis}>
          <span>{moment(utime * 1000).format("YYYY-MM-DD hh:mm:ss")}</span>
        </div>
      ),
    },
    {
      title: i18n.formatMessage({ id: "status" }),
      dataIndex: "status",
      width: 100,
      align: "center",
      render: (value) => {
        const status = AlarmStatus.find((item) => value === item.status);
        if (!status) return <>-</>;
        return (
          <div>
            <IconFont
              type={status.icon}
              size={100}
              style={{ color: status.color, marginRight: "10px" }}
            />
            <a>{status.label}</a>
          </div>
        );
      },
    },
    {
      title: i18n.formatMessage({ id: "operation" }),
      dataIndex: "operations",
      align: "center",
      width: 180,
      render: (_: any, record: AlarmType) => (
        <Space>
          {record.status !== 0 && (
            <>
              <Tooltip
                title={i18n.formatMessage({
                  id: `alarm.rules.switch.${
                    record.status === 1 ? "open" : "close"
                  }`,
                })}
              >
                <a onClick={() => doUpdateStatus(record)}>
                  <IconFont
                    type={record.status === 1 ? "icon-play" : "icon-suspended"}
                  />
                </a>
              </Tooltip>
              <Divider type="vertical" />
            </>
          )}

          <Tooltip
            title={i18n.formatMessage({
              id: "edit",
            })}
          >
            <a onClick={() => handleEdit(record)}>
              <EditOutlined />
            </a>
          </Tooltip>
          <Divider type="vertical" />
          <Tooltip title={i18n.formatMessage({ id: "alarm.rules.info.title" })}>
            <a onClick={() => handleInfo(record.id)}>
              <FileTextOutlined />
            </a>
          </Tooltip>
          <Divider type="vertical" />
          <Tooltip
            title={i18n.formatMessage({
              id: "delete",
            })}
          >
            <a onClick={() => handleDelete(record)}>
              <IconFont type={"icon-delete"} />
            </a>
          </Tooltip>
        </Space>
      ),
    },
  ];

  return (
    <div className={alarmStyles.tableMain}>
      <Table
        loading={doGetAlarms.loading}
        rowKey={"id"}
        size={"small"}
        columns={column}
        dataSource={alarmList}
        rowClassName={alarmStyles.tableWrapper}
        pagination={{
          responsive: true,
          showSizeChanger: true,
          size: "small",
          ...currentPagination,
          onChange: (page, pageSize) => {
            onChangePagination({
              ...currentPagination,
              current: page,
              pageSize,
            });
            doGetAlarms.run({ ...searchQuery, current: page, pageSize });
          },
        }}
      />
    </div>
  );
};
export default AlarmTable;
