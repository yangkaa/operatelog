<template>
  <div class="apps-container">
    <div class="titles">日志记录</div>
    <div style="margin-top: 20px">
      <el-form
        :inline="true"
        @submit.native.prevent
        :model="listQuery"
        size="small"
        label-width="300px"
      >
        <el-form-item class="searchStyle">
          <el-input
            v-model="listQuery.query"
            style="width: 500px;"
            size="small"
            placeholder="请输入操作描述"
            @keyup.enter.native="handleSearchList"
          >
            <el-button
              slot="append"
              icon="el-icon-search"
              @click="handleSearchList"
            ></el-button>
          </el-input>
        </el-form-item>
        <el-date-picker
          class="zslFormtime"
          size="small"
          type="datetimerange"
          range-separator="-"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="yyyy-MM-dd HH:mm:ss"
          value-format="yyyy-MM-dd HH:mm:ss"
          :picker-options="expireTimeOption"
          @change="handleChange"
          :default-time="['00:00:00', endTime]"
          v-model="searchItems"
        ></el-date-picker>
      </el-form>
    </div>

    <div class="table-container">
      <el-table
        v-loading="listLoading"
        :border="true"
        ref="homeBrandTable"
        :data="list"
        style="width: 100%;"
      >
        <el-table-column label="创建时间" width="250" align="center">
          <template slot-scope="scope">
            {{ scope.row.create_time ? timetrans(scope.row.create_time) : "-" }}
          </template>
        </el-table-column>

        <el-table-column
          prop="ip"
          label="IP"
          width="250"
          align="center"
        ></el-table-column>
        <el-table-column
          prop="log_level"
          label="日志等级"
          width="180"
          align="center"
        >
        </el-table-column>
        <el-table-column
          prop="log_type"
          label="日志类型"
          width="180"
          align="center"
        >
          <template slot-scope="scope">
            {{ scope.row.log_type ? handleLogLevel(scope.row.log_type) : "-" }}
          </template>
        </el-table-column>
        <el-table-column
          prop="op_name"
          label="操作名称"
          width="180"
          align="center"
        ></el-table-column>
        <el-table-column prop="op_desc" label="操作描述"> </el-table-column>
      </el-table>
    </div>

    <div class="pagination-container">
      <el-pagination
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        layout="total, sizes,prev, pager, next,jumper"
        :page-size="listQuery.pageSize"
        :page-sizes="[5, 10, 15]"
        :current-page.sync="listQuery.page"
        :total="total"
      ></el-pagination>
    </div>
  </div>
</template>
<script>
import axios from "axios";

const defaultListQuery = {
  page: 1,
  pageSize: 10,
  query: null,
  startTime: "",
  endTime: ""
};

export default {
  name: "EnterpriseManagement",
  data() {
    return {
      display_time: "",
      searchItems: [],
      listQuery: Object.assign({}, defaultListQuery),
      list: [],
      total: null,
      listLoading: true,
      endTime: this.moment(new Date())
        .locale("zh-cn")
        .format("HH:mm:ss"),
      expireTimeOption: {
        disabledDate(time) {
          return time.getTime() > Date.now();
        }
      }
    };
  },
  created() {
    this.getEnterpriseList();
  },
  methods: {
    handleChange() {
      if (this.searchItems && this.searchItems.length > 1) {
        this.listQuery.startTime = this.searchItems[0];
        this.listQuery.endTime = this.searchItems[1];
      } else {
        this.listQuery.startTime = "";
        this.listQuery.endTime = "";
      }
      this.listQuery.page = 1;
      this.getEnterpriseList();
    },
    handleLogLevel(type) {
      return type == "1" ? "审计日志" : type;
    },
    timetrans(dates) {
      if (dates && dates.indexOf("0001-01") > -1) {
        return "-";
      }
      const deter = new Date(dates).getTime();

      const date = new Date(deter); // 如果date为13位不需要乘1000
      const Y = date.getFullYear() + "-";
      const M =
        (date.getMonth() + 1 < 10
          ? "0" + (date.getMonth() + 1)
          : date.getMonth() + 1) + "-";
      const D =
        (date.getDate() < 10 ? "0" + date.getDate() : date.getDate()) + " ";
      const h =
        (date.getHours() < 10 ? "0" + date.getHours() : date.getHours()) + ":";
      const m =
        (date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes()) +
        ":";
      const s =
        date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();
      const times = Y + M + D + h + m + s;
      return times;
    },
    handleSizeChange(val) {
      this.listQuery.page = 1;
      this.listQuery.pageSize = val;
      this.getEnterpriseList();
    },
    handleCurrentChange(val) {
      this.listQuery.page = val;
      this.getEnterpriseList();
    },
    handleSearchList() {
      this.listQuery.page = 1;
      this.getEnterpriseList();
    },
    getEnterpriseList() {
      this.listLoading = true;
      axios({
        url: "/api/v1/mkyAuditLogs",
        method: "get",
        params: this.listQuery,
        headers: {
          Authorization: "test@token"
        }
      })
        .then(re => {
          if (re && re.data) {
            this.list = re.data.logs || [];
            this.total = re.data.totalCount || null;
          }
          this.listLoading = false;
        })
        .catch(() => {
          this.listLoading = false;
        });
    }
  }
};
</script>
<style scoped>
.zslFormtime {
  float: right;
}
.pagination-container {
  margin-top: 30px;
}
.searchStyle {
  float: left;
}
.apps-container {
  width: 90%;
  padding: 50px;
  margin: 20px auto 0;
  background: #fff;
}
.titles {
  font-size: 19px;
  height: 30px;
  line-height: 30px;
  margin-bottom: 40px;
  padding-bottom: 40px;
  border-bottom: 1px solid #eeeeee;
}
</style>
